package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

func ServeHTTP(builders []qg.GameBuilder, completeFn func(qg.Game), qgnPath, port string) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %v", l.Addr())

	raw, err := os.ReadFile(qgnPath)
	if err != nil {
		log.Fatal(err)
	}
	snapshot, err := qgn.Parse(string(raw))
	if err != nil {
		log.Fatal(err)
	}

	var builder qg.GameBuilder
	for _, b := range builders {
		if b.GetInformation().Key == snapshot.Tags[qgn.KeyTag] {
			builder = b
			break
		}
	}
	if builder == nil {
		log.Fatal(fmt.Errorf("no builder found for %s", snapshot.Tags[qgn.KeyTag]))
	}
	game, err := builder.Create(snapshot)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Handler:      NewGameServer(game, completeFn),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
