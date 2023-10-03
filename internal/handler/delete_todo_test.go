package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-weed-backend/internal/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestDeleteTodo(t *testing.T) {
	db, cleanup := setupTestDatabase()
	defer cleanup()

	Init(db)

	// ダミーデータの作成
	dummyTodo := model.Todo{Title: "Task to Delete", Completed: false}
	db.Create(&dummyTodo)

	req := httptest.NewRequest("DELETE", "/todos/delete?ID="+fmt.Sprint(dummyTodo.ID), nil)
	w := httptest.NewRecorder()

	DeleteTodo(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// DBからデータを取得し、データが存在しないことを確認
	var todo model.Todo
	if err := db.First(&todo, dummyTodo.ID).Error; err == nil {
		t.Errorf("Expected todo to be deleted, but was found in db")
	} else if !gorm.IsRecordNotFoundError(err) {
		t.Errorf("Unexpected error occurred: %v", err)
	}
}
