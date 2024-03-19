package quibbble_controller

import "k8s.io/client-go/kubernetes"

func Clean(clientset *kubernetes.Clientset) error {
	c := NewController(clientset)
	return c.clean()
}
