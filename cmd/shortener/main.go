package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ShortenedURL struct {
	ID        int    `json:"id"`
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
}

var items []Item
var shortenedURLs []ShortenedURL

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Сервер запущен на порту 8080...")
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		addItem(w, r)
	case "GET":
		getOriginalURL(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func addItem(w http.ResponseWriter, r *http.Request) {
	// Установка заголовка Content-Type для того, чтобы гарантировать, что данные интерпретируются как JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	newItem.ID = len(items) + 1
	items = append(items, newItem)

	// Проверяем, что строка URL передана корректно
	if newItem.Name == "" {
		http.Error(w, "Пустая строка URL", http.StatusBadRequest)
		return
	}

	// Создаем сокращенную ссылку
	shortenedURL := ShortenedURL{
		ID:        newItem.ID,
		Original:  newItem.Name,
		Shortened: "http://localhost:8080/" + strconv.Itoa(newItem.ID),
	}
	shortenedURLs = append(shortenedURLs, shortenedURL)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shortenedURL.Shortened)
}

func getOriginalURL(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		http.Error(w, "Некорректный идентификатор сокращенной ссылки", http.StatusBadRequest)
		return
	}

	for _, shortenedURL := range shortenedURLs {
		if shortenedURL.ID == id {
			// Устанавливаем заголовок Location для перенаправления
			w.Header().Set("Location", shortenedURL.Original)
			// Устанавливаем статус ответа на 307 Temporary Redirect
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}

	// Если сокращенный URL не найден, отправляем ошибку 400
	http.Error(w, fmt.Sprintf("Сокращенный URL с ID %v не найден", id), http.StatusBadRequest)
}
