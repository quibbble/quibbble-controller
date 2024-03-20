package watcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	qs "github.com/quibbble/quibbble-controller/internal/server"
	"github.com/quibbble/quibbble-controller/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const timeout = time.Second * 3

type Watcher struct {
	clientset *kubernetes.Clientset
}

func NewWatcher(clientset *kubernetes.Clientset) *Watcher {
	return &Watcher{
		clientset: clientset,
	}
}

func (w *Watcher) list() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := w.clientset.CoreV1().Pods(k8s.Namespace).List(ctx, metav1.ListOptions{
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

func (w *Watcher) active(key, id string) (bool, error) {
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

	if active.PlayerCount > 0 || active.LastUpdated.Add(30*time.Minute).After(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (w *Watcher) delete(key, id string) error {
	url := fmt.Sprintf("http://quibbble-controller.%s/delete?key=%s&id=%s", k8s.Namespace, key, id)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete")
	}
	return nil
}
