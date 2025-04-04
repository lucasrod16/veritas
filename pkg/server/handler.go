package server

import (
	"net/http"

	"github.com/lucasrod16/veritas/pkg/printer"
	"github.com/lucasrod16/veritas/pkg/scanner"
)

func scanReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userInput := r.URL.Query().Get("image")
	if userInput == "" {
		http.Error(w, "Missing 'image' query parameter", http.StatusBadRequest)
		return
	}

	cfg, err := scanner.Scan(userInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cdx := printer.NewCycloneDXPrinter()

	reportPayload, err := cdx.Print(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(reportPayload))
}

func scanDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userInput := r.URL.Query().Get("image")
	if userInput == "" {
		http.Error(w, "Missing 'image' query parameter", http.StatusBadRequest)
		return
	}

	cfg, err := scanner.Scan(userInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vd := printer.NewVulnDetailsPrinter()

	vulnDetails, err := vd.Print(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(vulnDetails))
}

func stripSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' {
			newPath := path[:len(path)-1]
			r.URL.Path = newPath
		}
		next.ServeHTTP(w, r)
	})
}
