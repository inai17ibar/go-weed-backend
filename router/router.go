package router

import (
	"encoding/json"
	"go-weed-backend/internal/handler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

type Config struct {
	ServerPort string     `json:"ServerPort"`
	CORS       CORSConfig `json:"CORS"`
}

type CORSConfig struct {
	AllowedOrigins []string `json:"AllowedOrigins"`
	AllowedMethods []string `json:"AllowedMethods"`
	AllowedHeaders []string `json:"AllowedHeaders"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	env := os.Getenv("APP_ENV")
	var configFilePath string
	if env == "production" {
		configFilePath = "config.production.json"
	} else {
		configFilePath = "config.local.json"
	}

	config, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins(config.CORS.AllowedOrigins),
		handlers.AllowedMethods(config.CORS.AllowedMethods),
		handlers.AllowedHeaders(config.CORS.AllowedHeaders),
	)

	mux.Handle("/todos", corsHandler(http.HandlerFunc(handler.GetTodos)))
	mux.Handle("/addTodo", corsHandler(http.HandlerFunc(handler.AddTodo)))
	mux.Handle("/todos/delete", corsHandler(http.HandlerFunc(handler.DeleteTodo)))
	mux.Handle("/todos/update", corsHandler(http.HandlerFunc(handler.UpdateTodo)))
	mux.Handle("/todosByDate", corsHandler(http.HandlerFunc(handler.GetTodosByDate)))
	mux.Handle("/commits", corsHandler(http.HandlerFunc(handler.GetCommits)))
	mux.Handle("/commitDataByDate", corsHandler(http.HandlerFunc(handler.AggregateCommitDataByDate)))

	return mux
}
