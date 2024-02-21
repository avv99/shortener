package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"shortener/internal/app/handlers"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	serverAddress := getEnv("SERVER_ADDRESS", ":8080")
	baseURL := getEnv("BASE_URL", "http://localhost:8080/")

	handlers.SetBaseURL(baseURL) // Установка baseURL для обработчиков

	r := chi.NewRouter()

	r.Get("/{id}", handlers.GetOriginalURL)
	r.Post("/", handlers.AddItem)
	r.Post("/api/shorten", handlers.APIShorten)

	fmt.Printf("Сервер запущен на адресе %s\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s\n", err)
	}
}
