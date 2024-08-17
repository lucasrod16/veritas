package veritashttp

import (
	"net/http"

	"github.com/anchore/grype/grype"
	"github.com/anchore/grype/grype/db"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Welcome to Veritas 🤠\n"))
		return
	}
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func scanHandler(dbCfg db.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _, closer, err := grype.LoadVulnerabilityDB(dbCfg, true)
			if err != nil {
				http.Error(w, "failed to load vulnerability database", http.StatusInternalServerError)
				return
			}
			defer closer.Close()
			w.Write([]byte("successfully loaded vulnerability database 🔐\n"))
			return
		}
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
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
