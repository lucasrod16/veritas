package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", routeHandler)

	fmt.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		if r.Method == http.MethodGet {
			w.Write([]byte("Welcome to Veritas 🤠\n"))
			return
		}
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	case "/test", "/test/":
		if r.Method == http.MethodGet {
			w.Write([]byte("hello world\n"))
			return
		}
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	default:
		http.NotFound(w, r)
	}
}
