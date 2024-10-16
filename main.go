package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasrod16/veritas/pkg/server"
)

//go:embed dashboard/*
//go:embed dashboard/images/*
var dashboard embed.FS

func main() {
	server.Dashboard = dashboard

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listening on %s\n", server.Addr)

	<-ctx.Done()
	stop()
	log.Println("Server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server gracefully shut down")
}
