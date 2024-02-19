package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
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

	fmt.Println("Сервер запущен на порту 8080...")
	http.ListenAndServe(":8080", r)
}
