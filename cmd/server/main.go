package main

import (
	"log"
	"os"

	"github.com/quibbble/quibbble-controller/games/tictactoe"
	qs "github.com/quibbble/quibbble-controller/internal/server"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	path := os.Getenv("QGN_PATH")
	if path == "" {
		path = "./qgn"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server starting...")
	defer log.Println("server closed")

	completeFn := func(qg.Game) {
		// todo store to complete table
	}

	qs.ServeHTTP([]qg.GameBuilder{tictactoe.Builder{}}, completeFn, path, port)
}
