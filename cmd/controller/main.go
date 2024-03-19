package main

import (
	"fmt"
	"log"
	"os"

	qc "github.com/quibbble/quibbble-controller/internal/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const (
	createMode = "CREATE"
	cleanMode  = "CLEAN"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mode := os.Getenv("MODE")
	if mode == "" {
		mode = createMode
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	switch mode {
	case createMode:
		qc.ServeHTTP(clientset, port)
	case cleanMode:
		if err := qc.Clean(clientset); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(fmt.Errorf("mode %s not valid", mode))
	}
}
