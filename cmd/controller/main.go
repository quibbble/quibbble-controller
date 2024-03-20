package main

import (
	"log"
	"os"

	qc "github.com/quibbble/quibbble-controller/internal/controller"
	crdb "github.com/quibbble/quibbble-controller/pkg/store/cockroachdb"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Config struct {
	Storage *crdb.Config         `yaml:"storage"`
	Server  *qc.GameServerConfig `yaml:"server"`
	Port    string               `yaml:"port"`
}

func main() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./config.yaml"
	}
	storagePassword := os.Getenv("STORAGE_PASSWORD")

	// creates the in-cluster config
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatal(err)
	}

	// read in configs
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	if err = yaml.Unmarshal(f, &config); err != nil {
		log.Fatal(err)
	}
	config.Storage.Password = storagePassword
	storage, err := crdb.NewClient(config.Storage)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("controller starting...")
	defer log.Println("controller closed")

	qc.ServeHTTP(clientset, storage, config.Server, config.Port)
}
