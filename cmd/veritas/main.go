package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/lucasrod16/veritas/pkg/veritashttp"
)

func main() {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt)

	err := veritashttp.StartServer(context.Background(), shutdownSignal)
	if err != nil {
		log.Fatal(err)
	}
}
