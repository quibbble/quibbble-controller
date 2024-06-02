package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/quibbble/quibbble-controller/games"
	qs "github.com/quibbble/quibbble-controller/internal/server"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	st "github.com/quibbble/quibbble-controller/pkg/store"
	crdb "github.com/quibbble/quibbble-controller/pkg/store/cockroachdb"
	"gopkg.in/yaml.v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Config struct {
	Storage *crdb.Config `yaml:"storage"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	qgnPath := os.Getenv("QGN_PATH")
	if qgnPath == "" {
		qgnPath = "./qgn"
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "quibbble"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	storagePassword := os.Getenv("STORAGE_PASSWORD")

	// parse qgn snapshot
	raw, err := os.ReadFile(qgnPath)
	if err != nil {
		log.Fatal(err)
	}
	snapshot, err := qgn.Parse(string(raw))
	if err != nil {
		log.Fatal(err)
	}

	// read in configs
	f, err := os.ReadFile(configPath)
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

	// create completeFn
	completeFn := func(snapshot *qgn.Snapshot) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		storage.StoreCompletedGame(ctx, &st.Game{
			Key:       snapshot.Tags[qgn.KeyTag],
			Snapshot:  snapshot,
			UpdatedAt: time.Now(),
		})
	}

	// retrieve builder
	var builder qg.GameBuilder
	for _, b := range games.Builders {
		if b.GetInformation().Key == snapshot.Tags[qgn.KeyTag] {
			builder = b
			break
		}
	}
	if builder == nil {
		log.Fatal(fmt.Errorf("no builder found for %s", snapshot.Tags[qgn.KeyTag]))
	}

	log.Println("server starting...")
	defer log.Println("server closed")

	qs.ServeHTTP(&qs.Params{
		Builder:       builder,
		Port:          port,
		CompleteFn:    completeFn,
		Snapshot:      snapshot,
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
	})
}
