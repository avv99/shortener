package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/config"
	"testing"
)

func TestAddItem(t *testing.T) {
	// Создание тестового запроса
	requestBody := []byte("http://example.com")
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Создание mock HTTP ResponseWriter
	w := httptest.NewRecorder()

	// Установка baseURL для тестового окружения
	baseURL := GetServerAddress()
	SetBaseURL(baseURL)
	cfg := config.InitConfig()
	// Вызов функции
	AddItem(w, req, cfg)

	// Проверка статуса кода
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Проверка корректности ответа
	expectedResponse := baseURL + "/1"
	if w.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expectedResponse)
	}

	// Проверка добавления элемента
	if len(shortenedURLs) != 1 {
		t.Errorf("handler didn't add item to shortenedURLs slice")
	}

	// Проверка данных добавленного элемента
	addedItem := shortenedURLs[0]
	if addedItem.ID != 1 {
		t.Errorf("added item has wrong ID: got %v want %v", addedItem.ID, 1)
	}
	if addedItem.Original != "http://example.com" {
		t.Errorf("added item has wrong Original URL: got %v want %v", addedItem.Original, "http://example.com")
	}
	if addedItem.Shortened != expectedResponse {
		t.Errorf("added item has wrong Shortened URL: got %v want %v", addedItem.Shortened, expectedResponse)
	}
}
