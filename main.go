package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/anchore/grype/grype"
	"github.com/anchore/grype/grype/db"
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
		_, _, closer, err := grype.LoadVulnerabilityDB(newGrypeDBCfg(), true)
		if err != nil {
			http.Error(w, formatHTTPError(err, "unable to load vulnerability database"), http.StatusInternalServerError)
			return
		}
		defer closer.Close()
		w.Write([]byte("successfully loaded vulnerability database 🔐\n"))
		return
	}
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func formatHTTPError(err error, message string) string {
	return fmt.Sprintf("%s: %s", message, err.Error())
}

func newGrypeDBCfg() db.Config {
	return db.Config{
		DBRootDir:  filepath.Join(xdg.CacheHome, "veritas", "db"),
		ListingURL: "https://toolbox-data.anchore.io/grype/databases/listing.json",
	}
}
