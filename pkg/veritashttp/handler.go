package veritashttp

import (
	"net/http"

	"github.com/lucasrod16/veritas/pkg/scanner"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Welcome to Veritas 🤠\n"))
		return
	}
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userInput := r.PathValue("userInput")
		reportPayload, err := scanner.Scan(userInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(reportPayload))
		return
	}
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
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
