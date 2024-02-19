package main

import (
	"fmt"
	"net/http"
	"time"
)

var urls map[string]string

func main() {
	urls = make(map[string]string)

	http.HandleFunc("/", handlePost)
	http.HandleFunc("/{id}", handleGet)

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	id := generateID()
	urls[id] = url

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Shortened URL: http://localhost:8080/%s", id)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[1:]
	url, ok := urls[id]
	if !ok {
		http.Error(w, "URL Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
