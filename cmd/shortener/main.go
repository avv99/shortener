package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

var (
	shortURLs map[string]string // Мапа для хранения сокращенных URL
	baseURL   = "http://localhost:8080"
)

func main() {
	shortURLs = make(map[string]string)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/redirect/", handleRedirect)

	fmt.Println("Server is running at", baseURL)
	http.ListenAndServe(":8080", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	longURL := r.Form.Get("url")
	if longURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortID := generateShortID()
	shortURL := baseURL + "/redirect/" + shortID
	shortURLs[shortID] = longURL

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, shortURL)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortID := r.URL.Path[len("/redirect/"):]
	longURL, ok := shortURLs[shortID]
	if !ok {
		http.Error(w, "Invalid short URL", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func generateShortID() string {
	return strconv.Itoa(rand.Intn(1000000)) // Генерация случайного числа в диапазоне [0, 1000000)
}
