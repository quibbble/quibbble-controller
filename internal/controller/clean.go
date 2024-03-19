package controller

import (
	"github.com/quibbble/quibbble-controller/internal/controller/k8s"
	"k8s.io/client-go/kubernetes"
)

func Clean(clientset *kubernetes.Clientset) error {
	c := NewController(clientset)
	names, err := c.list()
	if err != nil {
		return err
	}
	for _, name := range names {
		key, id := k8s.KeyID(name)
		active, err := c.active(key, id)
		if err != nil {
			return err
		}
		if !active {
			if err := c.delete(key, id); err != nil {
				return err
			}
		}
	}
	return nil
}
