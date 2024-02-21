package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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
var storageMutex sync.Mutex

func init() {
	loadDataFromDisk()
}

func loadDataFromDisk() {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	if filePath == "" {
		log.Println("Переменная окружения FILE_STORAGE_PATH не установлена. Используется хранение в памяти.")
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Ошибка чтения файла хранилища: %v. Используется хранение в памяти.\n", err)
		return
	}

	err = json.Unmarshal(data, &shortenedURLs)
	if err != nil {
		log.Printf("Ошибка декодирования данных из файла хранилища: %v. Используется хранение в памяти.\n", err)
		return
	}

	log.Println("Данные успешно загружены из файла хранилища.")
}

func saveDataToDisk() {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	if filePath == "" {
		log.Println("Переменная окружения FILE_STORAGE_PATH не установлена. Невозможно сохранить данные на диск.")
		return
	}

	data, err := json.Marshal(shortenedURLs)
	if err != nil {
		log.Printf("Ошибка кодирования данных для сохранения на диск: %v. Невозможно сохранить данные.\n", err)
		return
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("Ошибка записи данных на диск: %v. Невозможно сохранить данные.\n", err)
		return
	}

	log.Println("Данные успешно сохранены на диск.")
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	str, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	id := len(shortenedURLs) + 1
	shortenedURL := ShortenedURL{
		ID:        id,
		Original:  string(str),
		Shortened: fmt.Sprintf("%s/%d", getBaseURL(r), id),
	}
	shortenedURLs = append(shortenedURLs, shortenedURL)

	saveDataToDisk()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL.Shortened))
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

	storageMutex.Lock()
	defer storageMutex.Unlock()

	newItem.ID = len(shortenedURLs) + 1

	newItem.Shortened = fmt.Sprintf("%s/%d", getBaseURL(r), newItem.ID)

	shortenedURLs = append(shortenedURLs, newItem)

	saveDataToDisk()

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

func GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idnew, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Некорректный идентификатор сокращенной ссылки")
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

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

func getBaseURL(r *http.Request) string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return baseURL
}
