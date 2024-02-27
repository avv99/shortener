package config

import (
	"flag"
	"os"
)

type Config struct {
	PORT        string
	HOST        string
	BaseURL     string
	FilePath    string
	DSN         string
	TypeStorage string
}

func InitConfig() *Config {
	config := Config{}
	address := flag.String("a", ":8080", "HTTP server address")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL for shortened URLs")
	filePath := flag.String("f", "", "Path to file with shortened URLs")
	flag.Parse()

	// Задание значений по умолчанию для переменных окружения
	if os.Getenv("SERVER_ADDRESS") == "" {
		if *address == "" {
			config.PORT = ":8080"
		} else {
			config.PORT = *address
		}
		os.Setenv("SERVER_ADDRESS", config.PORT)
	}
	if os.Getenv("BASE_URL") == "" {
		if *baseURL == "" {
			config.BaseURL = "http://localhost/"
		} else {
			config.BaseURL = *baseURL
		}
		os.Setenv("BASE_URL", config.BaseURL)
	}

	if os.Getenv("FILE_STORAGE_PATH") == "" {
		if *filePath == "" {
			config.TypeStorage = "inMemory"
			config.FilePath = ""
		} else {
			config.TypeStorage = "inFile"
			config.FilePath = *filePath
		}
		os.Setenv("FILE_STORAGE_PATH", config.FilePath)
	}
	return &config
}
