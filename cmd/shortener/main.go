package main

import (
	"fmt"
	"net/http"
)

func main() {
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

	// Чтение строки URL из тела запроса
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Ваша логика для сокращения URL
	shortenedURL := shortenURL(url)

	// Отправка ответа с сокращенным URL
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, shortenedURL)
}

// Ваша логика для сокращения URL
func shortenURL(url string) string {
	// Ваш код для сокращения URL
	return "shortened-url"
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]
	fmt.Fprintf(w, "This is GET /%s endpoint", id)
}
