package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/quibbble/quibbble-controller/games/tictactoe"
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
	storagePassword := os.Getenv("STORAGE_PASSWORD")

	// create qgn snapshot
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

	completeFn := func(game qg.Game) {
		snapshot, err := game.GetSnapshotQGN()
		if err != nil {
			log.Println(err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		storage.StoreComplete(ctx, &st.Game{
			Key:       snapshot.Tags[qgn.KeyTag],
			Snapshot:  snapshot,
			UpdatedAt: time.Now(),
		})
	}

	log.Println("server starting...")
	defer log.Println("server closed")

	qs.ServeHTTP([]qg.GameBuilder{tictactoe.Builder{}}, completeFn, snapshot, port)
}
