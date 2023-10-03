package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-weed-backend/internal/model"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestUpdateTodo(t *testing.T) {
	db, cleanup := setupTestDatabase()
	defer cleanup()

	Init(db)

	// ダミーデータの作成
	dummyTodo := model.Todo{Title: "Task to Update", Completed: false}
	db.Create(&dummyTodo)

	// 変更内容
	updateData := struct {
		Title     string `json:"Title"`
		Completed bool   `json:"Completed"`
	}{
		Title:     "Updated Task",
		Completed: true,
	}
	body, _ := json.Marshal(updateData)

	// リクエストの作成
	req := httptest.NewRequest("PUT", "/todos/update?ID="+fmt.Sprint(dummyTodo.ID), bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// ハンドラの呼び出し
	UpdateTodo(w, req)

	// レスポンスの取得
	resp := w.Result()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// DBからデータを取得し、更新されているか確認
	var updatedTodo model.Todo
	db.First(&updatedTodo, dummyTodo.ID)

	if updatedTodo.Title != updateData.Title || updatedTodo.Completed != updateData.Completed {
		t.Errorf("Updated todo does not match expected values. Got: %+v", updatedTodo)
	}
}
