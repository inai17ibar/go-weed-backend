package router

import (
	"go-weed-backend/internal/handler"
	"net/http"

	"github.com/gorilla/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
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
