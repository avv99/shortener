package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strconv"
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

var items []Item
var shortenedURLs []ShortenedURL

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
		Shortened: "http://localhost:8080/" + strconv.Itoa(id),
	}
	shortenedURLs = append(shortenedURLs, shortenedURL)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL.Shortened))
}

func GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idnew, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Некорректный идентификатор сокращенной ссылки")
		return
	}

	for _, shortenedURL := range shortenedURLs {
		if shortenedURL.ID == idnew {
			// Устанавливаем заголовок Location для перенаправления
			w.Header().Set("Location", shortenedURL.Original)
			// Устанавливаем статус ответа на 307 Temporary Redirect
			http.Redirect(w, r, shortenedURL.Original, http.StatusTemporaryRedirect)

			return
		}
	}

	// Если сокращенный URL не найден, отправляем ошибку 400
	http.Error(w, fmt.Sprintf("Сокращенный URL с ID %v не найден", id), http.StatusBadRequest)
}

func APIShorten(w http.ResponseWriter, r *http.Request) {
	var newItem ShortenedURL

	type result struct {
		Result string `json:"result"`
	}

	var resultW result

	// Декодирование JSON из тела запроса в переменную newItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерация ID для нового элемента
	newItem.ID = len(shortenedURLs) + 1

	// Формирование сокращенной ссылки
	//aeaeaeea
	newItem.Shortened = "http://localhost:8080/" + strconv.Itoa(newItem.ID)

	// Добавление нового элемента в массив
	shortenedURLs = append(shortenedURLs, newItem)

	resultW.Result = newItem.Shortened

	// Кодирование сокращенной ссылки в JSON
	ResponseJSON, err := json.Marshal(resultW)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа в виде JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(ResponseJSON)
}
