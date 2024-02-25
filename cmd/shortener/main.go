package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"shortener/internal/app/handlers"
)

func main() {
	// Объявление флагов
	serverAddress := flag.String("a", "localhost:8080", "адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	fileStoragePath := flag.String("f", "", "путь до файла с сокращёнными URL")
	flag.Parse()

	// Установка переменных окружения, если они не установлены через флаги
	if *fileStoragePath != "" {
		if err := os.Setenv("FILE_STORAGE_PATH", *fileStoragePath); err != nil {
			log.Fatalf("Не удалось установить переменную окружения FILE_STORAGE_PATH: %v\n", err)
		}
	}

	// Установка базового URL в обработчиках
	handlers.SetBaseURL(*baseURL)

	// Загрузка данных из файла (если есть)
	handlers.LoadDataFromDisk()

	r := chi.NewRouter()
	r.Get("/{id}", handlers.GetOriginalURL)
	r.Post("/", handlers.AddItem)
	r.Post("/api/shorten", handlers.APIShorten)

	// Формирование адреса сервера
	serverAddr := *serverAddress
	if serverAddr == "" {
		serverAddr = "localhost:8080"
	}

	fmt.Printf("Сервер запущен на адресе %s...\n", serverAddr)

	// Запуск HTTP-сервера
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v\n", err)
	}
}
