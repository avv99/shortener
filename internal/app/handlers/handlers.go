package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	items         []Item
	shortenedURLs []ShortenedURL
	baseURL       string // Глобальная переменная для хранения базового URL
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ShortenedURL struct {
	ID        int    `json:"id"`
	Original  string `json:"url"`
	Shortened string `json:"shortened"`
}

// SetBaseURL устанавливает базовый URL для формирования сокращенных ссылок
func SetBaseURL(url string) {
	baseURL = url
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	str, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	id := len(shortenedURLs) + 1
	shortenedURL := ShortenedURL{
		ID:        id,
		Original:  string(str),
		Shortened: fmt.Sprintf("%s%s", baseURL, strconv.Itoa(id)), // Использование baseURL
	}
	shortenedURLs = append(shortenedURLs, shortenedURL)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL.Shortened))
}

func GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	// Тело функции остается без изменений
}

func APIShorten(w http.ResponseWriter, r *http.Request) {
	var newItem ShortenedURL

	type result struct {
		Result string `json:"result"`
	}

	var resultW result

	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newItem.ID = len(shortenedURLs) + 1
	newItem.Shortened = fmt.Sprintf("%s%s", baseURL, strconv.Itoa(newItem.ID)) // Использование baseURL

	shortenedURLs = append(shortenedURLs, newItem)

	resultW.Result = newItem.Shortened

	ResponseJSON, err := json.Marshal(resultW)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(ResponseJSON)
}
