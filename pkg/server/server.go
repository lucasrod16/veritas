package server

import (
	"net/http"

	"golang.org/x/sync/errgroup"
)

func StartServer() (*http.Server, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", rootHandler)
	mux.HandleFunc("/scan", scanHandler)
	srv := &http.Server{Addr: ":8080", Handler: stripSlashes(mux)}

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
