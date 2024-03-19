package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"k8s.io/client-go/kubernetes"
)

func ServeHTTP(clientset *kubernetes.Clientset, port string) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %v", l.Addr())

	s := &http.Server{
		Handler:      NewController(clientset),
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
