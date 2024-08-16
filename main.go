package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", rootHandler)
	mux.HandleFunc("/test", testHandler)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", stripSlashes(mux)))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Welcome to Veritas 🤠\n"))
		return
	}
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("hello world\n"))
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
