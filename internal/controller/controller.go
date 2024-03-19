package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/quibbble/quibbble-controller/internal/controller/k8s"
	qs "github.com/quibbble/quibbble-controller/internal/server"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	timeout = time.Second * 3
)

type Controller struct {
	clientset *kubernetes.Clientset

	serveMux http.ServeMux
}

func NewController(clientset *kubernetes.Clientset) *Controller {
	c := &Controller{
		clientset: clientset,
	}
	c.serveMux.HandleFunc("/create", c.createHandler)
	c.serveMux.HandleFunc("/health", healthHandler)
	return c
}

func (c *Controller) find(key, id string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := c.clientset.CoreV1().Pods(k8s.Namespace).Get(ctx, k8s.Name(key, id), metav1.GetOptions{})
	return err == nil
}

func (c *Controller) create(snapshot *qgn.Snapshot) error {
	key := snapshot.Tags[qgn.KeyTag]
	id := snapshot.Tags[qgn.IDTag]

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := c.clientset.CoreV1().ConfigMaps(k8s.Namespace).Create(ctx, k8s.CreateConfigMap(snapshot), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.CoreV1().Pods(k8s.Namespace).Create(ctx, k8s.CreatePod(key, id), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.CoreV1().Services(k8s.Namespace).Create(ctx, k8s.CreateService(key, id), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.NetworkingV1().Ingresses(k8s.Namespace).Create(ctx, k8s.CreateIngress(key, id), metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

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

func (c *Controller) list() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := c.clientset.CoreV1().Pods(k8s.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", k8s.Component, k8s.GameComponent),
	})
	if err != nil {
		return nil, err
	}
	pods := make([]string, 0)
	for _, it := range l.Items {
		pods = append(pods, it.Name)
	}
	return pods, nil
}

func (c *Controller) active(key, id string) (bool, error) {
	url := fmt.Sprintf("http://%s.%s/active", k8s.Name(key, id), k8s.Namespace)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var active qs.Active
	if err := json.Unmarshal(body, &active); err != nil {
		return false, err
	}

	if active.Players > 0 || active.LastUpdated.Add(30*time.Minute).After(time.Now()) {
		return true, nil
	}
	return false, nil
}
