package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/lucasrod16/veritas/pkg/server"
)

//go:embed dashboard/*
//go:embed dashboard/images/*
var dashboard embed.FS

func main() {
	server.Dashboard = dashboard

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	server, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server listening on %s\n", server.Addr)

	<-shutdownSignal

	fmt.Println("Server shutting down...")

	err = server.Shutdown(ctx)
	cancel()
	if err != nil {
		log.Fatal()
	}
}
