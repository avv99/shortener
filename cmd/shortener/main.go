package main

import (
	"fmt"
	"net/http"
	"time"
)

// Мапа для хранения сокращенных URL
var urlMap map[string]string

func main() {
	urlMap = make(map[string]string)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/{id}", handleID)
	http.ListenAndServe(":8080", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "Bad Request: URL is required", http.StatusBadRequest)
		return
	}

	// Генерация уникального ID для сокращенного URL
	id := generateID()

	// Добавление сокращенного URL в мапу
	urlMap[id] = url

	// Формирование сокращенного URL
	shortURL := fmt.Sprintf("http://localhost:8080/%s", id)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, shortURL)
}

func handleID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/"):]
	if id == "" {
		http.Error(w, "Bad Request: ID is required", http.StatusBadRequest)
		return
	}

	// Поиск оригинального URL по ID в мапе
	originalURL, ok := urlMap[id]
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

// Генерация уникального ID
func generateID() string {
	// Здесь может быть логика генерации уникального ID, например, случайная строка
	// В данном примере просто используется текущее время в качестве ID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
