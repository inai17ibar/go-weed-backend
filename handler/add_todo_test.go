package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-weed-backend/model"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestAddTodo(t *testing.T) {
	// テスト用のデータベースをセットアップ
	db, cleanup := setupTestDatabase()
	defer cleanup()

	// ハンドラの初期化
	Init(db)

	// テスト用のHTTPリクエストを作成
	todo := model.Todo{Title: "Task 1", Completed: false}
	body, _ := json.Marshal(todo)
	req := httptest.NewRequest("POST", "/addTodo", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// ハンドラ関数を呼び出し
	AddTodo(w, req)

	// HTTPレスポンスを取得
	resp := w.Result()

	// レスポンスのステータスコードを確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// DBからデータを取得し、期待するデータが存在するか確認
	var todos []model.Todo
	db.Find(&todos)

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo to be added, got %d", len(todos))
	} else if todos[0].Title != "Task 1" || todos[0].Completed {
		t.Errorf("Added todo does not match expected values")
	}
}
