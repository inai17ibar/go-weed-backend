package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestGetTodos(t *testing.T) {
	// テスト用のデータベースをセットアップ
	db, cleanup := setupTestDatabase()
	defer cleanup()

	// テスト用のダミーTODOデータを挿入
	createDummyTodoData(db)

	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()

	// ハンドラ関数を呼び出し
	getTodos(w, req)

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
	var todos []Todo
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

// テスト用のデータベースセットアップ
func setupTestDatabase() (*gorm.DB, func()) {
	var err error
	db, err = gorm.Open("sqlite3", "test.db") // テスト用のSQLiteデータベースを使用
	if err != nil {
		log.Fatal(err)
	}

	// マイグレーションを実行してテーブルを作成
	db.AutoMigrate(&Todo{})

	// クリーンアップ用の関数を返す
	cleanup := func() {
		db.Close()
		if err := os.Remove("test.db"); err != nil {
			log.Printf("Failed to remove test database file: %v", err)
		}
	}

	return db, cleanup
}

// ダミーのTODOデータを作成しデータベースに挿入
func createDummyTodoData(db *gorm.DB) {
	dummyTodos := []Todo{
		{Title: "Task 1", Completed: false, Created_Date: "2023-09-01"},
		{Title: "Task 2", Completed: true, Created_Date: "2023-09-02"},
		{Title: "Task 3", Completed: false, Created_Date: "2023-09-03"},
	}

	for _, todo := range dummyTodos {
		db.Create(&todo)
	}
}
