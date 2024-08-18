package server

import (
	"net/http"

	"golang.org/x/sync/errgroup"
)

func StartServer(dashboardPath string) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(dashboardPath)))
	mux.HandleFunc("/scan/report", scanReportHandler)
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
