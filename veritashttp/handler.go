package veritashttp

import (
	"net/http"
	"strconv"

	"github.com/lucasrod16/veritas/scanner"
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
		matchCount, err := scanner.Scan(userInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matchStr := strconv.Itoa(matchCount) + "\n"
		w.Write([]byte(matchStr))
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
