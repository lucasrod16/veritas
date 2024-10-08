package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lucasrod16/veritas/pkg/server"
)

//go:embed dashboard/*
//go:embed dashboard/images/*
var dashboard embed.FS

func main() {
	server.Dashboard = dashboard

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listening on %s\n", server.Addr)

	<-shutdown

	log.Println("Server shutting down...")

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server gracefully shut down")
}
