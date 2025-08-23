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
	namespace string
	clientset *kubernetes.Clientset
}

func NewWatcher(namespace string, clientset *kubernetes.Clientset) *Watcher {
	return &Watcher{
		clientset: clientset,
	}
}

func (w *Watcher) list() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := w.clientset.CoreV1().Pods(w.namespace).List(ctx, metav1.ListOptions{
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
	url := fmt.Sprintf("http://%s.%s/activity", k8s.Name(key, id), w.namespace)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var active qs.Activity
	if err := json.Unmarshal(body, &active); err != nil {
		return false, err
	}

	// If the game has been inactive for an hour, regardless of player count, then return false
	if active.LastUpdated.Add(time.Hour).Before(time.Now()) {
		return false, nil
	}

	// If the game has 1 or more players or the last active was less than 15 minutes ago then return true
	if active.PlayerCount > 0 || active.LastUpdated.Add(15*time.Minute).After(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (w *Watcher) delete(key, id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://quibbble-controller.%s/game?key=%s&id=%s", w.namespace, key, id), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete")
	}
	return nil
}
