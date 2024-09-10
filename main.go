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

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	server, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listening on %s\n", server.Addr)

	<-shutdownSignal

	log.Println("Server shutting down...")

	err = server.Shutdown(ctx)
	cancel()
	if err != nil {
		log.Fatal()
	}

	log.Println("Server gracefully shut down")
}
