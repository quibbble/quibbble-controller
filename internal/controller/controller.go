package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	"github.com/quibbble/quibbble-controller/pkg/k8s"
	st "github.com/quibbble/quibbble-controller/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	timeout          = time.Second * 3
	maxLiveGameCount = 50 // the maximum number of concurrent live games that the controller will support.
)

type Controller struct {
	// clientset provides connection to the k8s cluster.
	clientset *kubernetes.Clientset

	// storage provides connection to the game store.
	storage st.GameStore

	// config for server creation.
	config *ServerConfig

	// mux handles http server handling.
	mux http.ServeMux
}

func NewController(config *ServerConfig, clientset *kubernetes.Clientset, storage st.GameStore) *Controller {
	c := &Controller{
		clientset: clientset,
		storage:   storage,
		config:    config,
	}
	c.mux.HandleFunc("POST /game", c.createHandler)
	c.mux.HandleFunc("PUT /game", c.loadHandler)
	c.mux.HandleFunc("DELETE /game", c.deleteHandler)
	c.mux.HandleFunc("GET /game/activity", c.activityHandler)
	c.mux.HandleFunc("GET /health", healthHandler)
	return c
}

// find checks to see if a game is currently live (game server up and running).
func (c *Controller) find(key, id string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := c.clientset.CoreV1().Pods(c.config.Namespace).Get(ctx, k8s.Name(key, id), metav1.GetOptions{})
	return err == nil
}

// create creates a new game server with the given snapshot.
func (c *Controller) create(snapshot *qgn.Snapshot) error {
	key := snapshot.Tags[qgn.KeyTag]
	id := snapshot.Tags[qgn.IDTag]

	port := rand.Int31n(65536)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := c.clientset.CoreV1().ConfigMaps(c.config.Namespace).Create(ctx, k8s.CreateConfigMap(snapshot), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.CoreV1().Pods(c.config.Namespace).Create(ctx, k8s.CreatePod(c.config.FullName, key, id, port, c.storage.Enabled(), &c.config.Pod), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.CoreV1().Services(c.config.Namespace).Create(ctx, k8s.CreateService(key, id), metav1.CreateOptions{}); err != nil {
		return err
	}
	if c.config.Ingress.Enabled {
		if _, err := c.clientset.NetworkingV1().Ingresses(c.config.Namespace).Create(ctx, k8s.CreateIngress(key, id, &c.config.Ingress), metav1.CreateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

// delete removes a game server with the given key and id.
func (c *Controller) delete(key, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var errList []error
	if err := c.clientset.CoreV1().ConfigMaps(c.config.Namespace).Delete(ctx, k8s.Name(key, id), metav1.DeleteOptions{}); err != nil {
		errList = append(errList, err)
	}
	if err := c.clientset.CoreV1().Pods(c.config.Namespace).Delete(ctx, k8s.Name(key, id), metav1.DeleteOptions{}); err != nil {
		errList = append(errList, err)
	}
	if err := c.clientset.CoreV1().Services(c.config.Namespace).Delete(ctx, k8s.Name(key, id), metav1.DeleteOptions{}); err != nil {
		errList = append(errList, err)
	}
	if c.config.Ingress.Enabled {
		if err := c.clientset.NetworkingV1().Ingresses(c.config.Namespace).Delete(ctx, k8s.Name(key, id), metav1.DeleteOptions{}); err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) > 0 {
		return errors.Join(errList...)
	}
	return nil
}

// lookup searches long term game storage to see if a game with key and id exists.
func (c *Controller) lookup(key, id string) (*qgn.Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	g, err := c.storage.GetActiveGame(ctx, key, id)
	if err != nil {
		return nil, err
	}
	return g.Snapshot, nil
}

// store saves a live game with key and id to long term game storage.
func (c *Controller) store(key, id string) error {
	url := fmt.Sprintf("http://%s.%s/snapshot?format=qgn", k8s.Name(key, id), c.config.Namespace)
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
	return c.storage.StoreActiveGame(ctx, &st.Game{
		Key:       key,
		ID:        id,
		Snapshot:  snapshot,
		UpdatedAt: time.Now(),
	})
}

func (c *Controller) increment(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.storage.IncrementGameCount(ctx, key)
}
