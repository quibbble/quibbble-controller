package main

import (
	"log"

	"github.com/quibbble/quibbble-controller/internal/watcher"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("watcher starting...")
	defer log.Println("watcher closed")

	if err := watcher.Clean(clientset); err != nil {
		log.Fatal(err)
	}
}
