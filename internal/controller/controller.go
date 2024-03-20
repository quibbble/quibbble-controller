package controller

import (
	"context"
	"net/http"
	"time"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	"github.com/quibbble/quibbble-controller/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const timeout = time.Second * 3

type Controller struct {
	clientset *kubernetes.Clientset

	serveMux http.ServeMux
}

func NewController(clientset *kubernetes.Clientset) *Controller {
	c := &Controller{
		clientset: clientset,
	}
	c.serveMux.HandleFunc("/create", c.createHandler)
	c.serveMux.HandleFunc("/delete", c.deleteHandler)
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
