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

type Params struct {
	// Builder that builds the game instance
	Builder qg.GameBuilder

	// Snapshot of the game to build
	Snapshot *qgn.Snapshot

	// Port opened on the server
	Port string

	// CompleteFn called on game end
	CompleteFn func(qg.Game)
}

func ServeHTTP(p *Params) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", p.Port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %v", l.Addr())

	game, err := p.Builder.Create(p.Snapshot)
	if err != nil {
		log.Fatal(err)
	}

	// retrieve server tags
	id, ok := p.Snapshot.Tags[qgn.IDTag]
	if !ok {
		log.Fatal(fmt.Errorf("missing id tag"))
	}
	kind, ok := p.Snapshot.Tags[qgn.KindTag]
	if !ok {
		log.Fatal(fmt.Errorf("missing kind tag"))
	}
	s := &http.Server{
		Handler:      NewGameServer(game, id, kind, p.CompleteFn),
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
