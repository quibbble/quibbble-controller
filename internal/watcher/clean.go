package watcher

import (
	"errors"

	"github.com/quibbble/quibbble-controller/pkg/k8s"
	"k8s.io/client-go/kubernetes"
)

func Clean(namespace string, clientset *kubernetes.Clientset) error {
	w := NewWatcher(namespace, clientset)
	names, err := w.list()
	if err != nil {
		return err
	}
	var errList []error
	for _, name := range names {
		key, id := k8s.KeyID(name)
		active, err := w.active(key, id)
		if err != nil {
			errList = append(errList, err)
			active = false
		}
		if !active {
			if err := w.delete(key, id); err != nil {
				errList = append(errList, err)
			}
		}
	}
	if len(errList) > 0 {
		return errors.Join(errList...)
	}
	return nil
}
