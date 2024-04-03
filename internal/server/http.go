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

	// Authenticate validates protected endpoints
	Authenticate func(http.Handler) http.Handler
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
	typ, err := p.Snapshot.Tags.Type()
	if err != nil {
		log.Fatal(err)
	}
	players, err := p.Snapshot.Tags.Players()
	if err != nil {
		log.Fatal(err)
	}
	id, ok := p.Snapshot.Tags[qgn.IDTag]
	if !ok {
		log.Fatal(fmt.Errorf("missing id tag"))
	}

	s := &http.Server{
		Handler:      NewGameServer(game, id, typ, players, p.CompleteFn, p.Authenticate),
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
