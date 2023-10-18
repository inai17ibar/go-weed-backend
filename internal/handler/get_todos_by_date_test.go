package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-weed-backend/internal/model"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestGetTodosByDate(t *testing.T) {
	db, cleanup := setupTestDatabase()
	defer cleanup()

	InitForTest(db)

	date := "2023-09-01"
	dummyTodo := model.Todo{Title: "Task 1", Completed: false, Created_Date: date}
	db.Create(&dummyTodo)

	req := httptest.NewRequest("GET", "/todosByDate?Created_date="+date, nil)
	w := httptest.NewRecorder()

	GetTodosByDate(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var todos []model.Todo
	json.NewDecoder(resp.Body).Decode(&todos)

	if len(todos) != 1 || todos[0].Title != "Task 1" {
		t.Errorf("Expected 1 todo with title 'Task 1', got %+v", todos)
	}
}
