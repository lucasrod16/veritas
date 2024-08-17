package veritashttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lucasrod16/veritas/vdb"
)

func StartServer(ctx context.Context, shutdownSignal <-chan os.Signal) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", rootHandler)
	mux.HandleFunc("/scan", scanHandler(vdb.NewGrypeDBCfg()))

	srv := &http.Server{Addr: ":8080", Handler: stripSlashes(mux)}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fmt.Printf("Server listening on %s\n", srv.Addr)

	<-shutdownSignal

	fmt.Println("Server shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
