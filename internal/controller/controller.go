package controller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	"github.com/quibbble/quibbble-controller/pkg/k8s"
	st "github.com/quibbble/quibbble-controller/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const timeout = time.Second * 3

type Controller struct {
	// clientset provides connection to the k8s cluster.
	clientset *kubernetes.Clientset

	// storage provides connection to the game store.
	storage st.GameStore

	// config for game server creation.
	config *GameServerConfig

	// serveMux handles http server handling.
	mux *chi.Mux
}

func NewController(config *GameServerConfig, clientset *kubernetes.Clientset, storage st.GameStore, allowedOrigins []string) *Controller {
	c := &Controller{
		clientset: clientset,
		storage:   storage,
		config:    config,
		mux:       chi.NewRouter(),
	}
	c.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{http.MethodHead, http.MethodPost, http.MethodGet, http.MethodDelete},
		AllowCredentials: true,
	}))
	c.mux.Post("/create", c.createHandler)
	c.mux.Delete("/delete", c.deleteHandler)
	c.mux.Get("/stats", c.statsHandler)
	c.mux.Get("/health", healthHandler)
	return c
}

// find checks to see if a game is currently live (game server up and running).
func (c *Controller) find(key, id string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := c.clientset.CoreV1().Pods(k8s.Namespace).Get(ctx, k8s.Name(key, id), metav1.GetOptions{})
	return err == nil
}

// create creates a new game server with the given snapshot.
func (c *Controller) create(snapshot *qgn.Snapshot) error {
	key := snapshot.Tags[qgn.KeyTag]
	id := snapshot.Tags[qgn.IDTag]

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := c.clientset.CoreV1().ConfigMaps(k8s.Namespace).Create(ctx, k8s.CreateConfigMap(snapshot), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.CoreV1().Pods(k8s.Namespace).Create(ctx, k8s.CreatePod(key, id, fmt.Sprintf("%s:%s", c.config.Image.Repository, c.config.Image.Tag), c.config.Image.PullPolicy), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.CoreV1().Services(k8s.Namespace).Create(ctx, k8s.CreateService(key, id), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.NetworkingV1().Ingresses(k8s.Namespace).Create(ctx, k8s.CreateIngress(c.config.Host, key, id), metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

// delete removes a game server with the given key and id.
func (c *Controller) delete(key, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := c.clientset.CoreV1().ConfigMaps(k8s.Namespace).Delete(ctx, k8s.Name(key, id), metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err := c.clientset.CoreV1().Pods(k8s.Namespace).Delete(ctx, k8s.Name(key, id), metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err := c.clientset.CoreV1().Services(k8s.Namespace).Delete(ctx, k8s.Name(key, id), metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err := c.clientset.NetworkingV1().Ingresses(k8s.Namespace).Delete(ctx, k8s.Name(key, id), metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

// lookup searches long term game storage to see if a game with key and id exists.
func (c *Controller) lookup(key, id string) (*qgn.Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	g, err := c.storage.GetGame(ctx, key, id)
	if err != nil {
		return nil, err
	}
	return g.Snapshot, nil
}

// store saves a live game with key and id to long term game storage.
func (c *Controller) store(key, id string) error {
	url := fmt.Sprintf("http://%s.%s/snapshot?format=qgn", k8s.Name(key, id), k8s.Namespace)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	snapshot, err := qgn.Parse(string(body))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.storage.StoreActive(ctx, &st.Game{
		Key:       key,
		ID:        id,
		Snapshot:  snapshot,
		UpdatedAt: time.Now(),
	})
}
