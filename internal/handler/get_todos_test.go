package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-weed-backend/internal/model"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestGetTodos(t *testing.T) {
	// テスト用のデータベースをセットアップ
	client, cleanup := setupTestDatabase()
	defer cleanup()

	// テスト用のダミーTODOデータを挿入
	createDummyTodoData(client)

	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()

	// ハンドラ関数を呼び出し
	GetTodos(w, req) // この関数もMongoDB対応に変更する必要があります

	// HTTPレスポンスを取得
	resp := w.Result()

	// レスポンスのステータスコードを確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// レスポンスのContent-Typeを確認
	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Expected Content-Type: %s, got %s", expectedContentType, actualContentType)
	}

	// レスポンスボディをJSONデコード
	var todos []model.Todo
	err := json.NewDecoder(resp.Body).Decode(&todos)
	if err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}

	// ダミーTODOデータが正しく取得できたか確認
	expectedTodoCount := 3 // ダミーデータが3つ挿入されていると仮定
	if len(todos) != expectedTodoCount {
		t.Errorf("Expected %d TODO items, got %d", expectedTodoCount, len(todos))
	}

	// ここでtodosを期待するデータと比較して、JSONレスポンスの内容を検証することができます

	// クリーンアップ
	resp.Body.Close()
}
