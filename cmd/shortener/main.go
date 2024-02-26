package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"shortener/internal/app/handlers"
)

func main() {
	// Определение флагов командной строки
	address := flag.String("a", ":8080", "HTTP server address")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL for shortened URLs")
	filePath := flag.String("f", "", "Path to file with shortened URLs")
	flag.Parse()

	// Задание значений по умолчанию для переменных окружения
	if os.Getenv("SERVER_ADDRESS") == "" {
		os.Setenv("SERVER_ADDRESS", *address)
	}
	if os.Getenv("BASE_URL") == "" {
		os.Setenv("BASE_URL", *baseURL)
	}
	if os.Getenv("FILE_STORAGE_PATH") == "" {
		os.Setenv("FILE_STORAGE_PATH", *filePath)
	}

	r := chi.NewRouter()
	r.Get("/{id}", handlers.GetOriginalURL)
	r.Post("/", handlers.AddItem)
	r.Post("/api/shorten", handlers.APIShorten)

	serverAddress := os.Getenv("SERVER_ADDRESS")
	fmt.Printf("Сервер запущен на адресе %s...\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v\n", err)
	}
}
