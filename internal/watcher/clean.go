package watcher

import (
	"github.com/quibbble/quibbble-controller/pkg/k8s"
	"k8s.io/client-go/kubernetes"
)

func Clean(clientset *kubernetes.Clientset) error {
	w := NewWatcher(clientset)
	names, err := w.list()
	if err != nil {
		return err
	}
	for _, name := range names {
		key, id := k8s.KeyID(name)
		active, err := w.active(key, id)
		if err != nil {
			return err
		}
		if !active {
			if err := w.delete(key, id); err != nil {
				return err
			}
		}
	}
	return nil
}
