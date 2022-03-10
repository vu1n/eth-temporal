package main

import (
	"eth-temporal/app"
	"log"
	"net/http"
	"os"
	"os/signal"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := app.NewClient(client.Options{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer c.Close()

	server := &http.Server{
		Handler: Router(c),
		Addr:    "0.0.0.0:8081",
	}

	errCh := make(chan error, 1)
	go func() { errCh <- server.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	select {
	case <-sigCh:
		server.Close()
	case err = <-errCh:
		log.Fatalf("error: %v", err)
	}
}
