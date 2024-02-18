package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type URLShortener struct {
	urlMap    map[int]string
	idCounter int
	mutex     sync.Mutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urlMap:    make(map[int]string),
		idCounter: 1,
		mutex:     sync.Mutex{},
	}
}

func (s *URLShortener) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var urlData struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&urlData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	shortURL := strconv.Itoa(s.idCounter)
	s.urlMap[s.idCounter] = urlData.URL
	s.idCounter++

	response := map[string]string{"shortened_url": fmt.Sprintf("http://localhost:8080/%s", shortURL)}
	responseJSON, _ := json.Marshal(response)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (s *URLShortener) RedirectOriginalURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if idInt, err := strconv.Atoi(id); err == nil {
		if originalURL, ok := s.urlMap[idInt]; ok {
			http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
			return
		}
	}
	http.Error(w, "URL not found", http.StatusBadRequest)
}

func main() {
	shortener := NewURLShortener()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			shortener.ShortenURL(w, r)
		case http.MethodGet:
			shortener.RedirectOriginalURL(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
