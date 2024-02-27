package main

import (
	"fmt"
	"log"
	"net/http"
	"shortener/internal/app/config"
	"shortener/internal/app/middleware"

	"github.com/go-chi/chi/v5"
	"shortener/internal/app/handlers"
)

func main() {
	// Определение флагов командной строки
	cfg := config.InitConfig()

	r := chi.NewRouter()
	//r.Get("/{id}", handlers.GetOriginalURL)
	//r.Post("/", handlers.AddItem)
	//r.Post("/api/shorten", handlers.APIShorten)
	r.Use(middleware.GzipMiddleware)
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.APIShorten(w, r, cfg)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOriginalURL(w, r)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddItem(w, r, cfg)
	})

	//serverAddress := os.Getenv("SERVER_ADDRESS")
	fmt.Printf("Сервер запущен на адресе %s...\n", cfg.PORT)
	if err := http.ListenAndServe(cfg.PORT, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v\n", err)
	}
}
