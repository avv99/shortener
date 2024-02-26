package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os" // Импортируем пакет для работы с переменными окружения1112
	"shortener/internal/app/handlers"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOriginalURL(w, r)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddItem(w, r)
	})

	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.APIShorten(w, r)
	})

	// Получаем адрес сервера из переменной окружения111
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		// Если переменная не установлена, используем порт по умолчанию
		serverAddress = ":8080"
	}

	fmt.Printf("Сервер запущен на адресе %s...\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
