package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"golang.org/x/sync/errgroup"
)

var Dashboard embed.FS

func Start() (*http.Server, error) {
	fs, err := fs.Sub(Dashboard, "dashboard")
	if err != nil {
		return nil, fmt.Errorf("failed to load dashboard assets: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(fs)))
	mux.HandleFunc("/scan/report", scanReportHandler)
	mux.HandleFunc("/scan/details", scanDetailsHandler)

	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: stripSlashes(mux)}

	var g errgroup.Group
	g.Go(func() error {
		err := srv.ListenAndServe()
		if err != nil {
			return err
		}
		return nil
	})

	return srv, nil
}
