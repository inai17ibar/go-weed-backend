package handler

import (
	"encoding/json"
	"go-weed-backend/model"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

func TestGetCommits(t *testing.T) {
	// テスト用のデータベースをセットアップ
	var cleanup func()
	db, cleanup = setupTestCommitsDatabase()
	defer cleanup()

	// ハンドラの初期化
	Init(db) // 追加

	// テスト用のダミーTODOデータを挿入
	createDummyCommitsData()

	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest("GET", "/commits", nil)
	w := httptest.NewRecorder()

	// ハンドラ関数を呼び出し
	GetCommits(w, req)

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
	var myCommits []model.MyCommit
	err := json.NewDecoder(resp.Body).Decode(&myCommits)
	if err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}

	// ダミーTODOデータが正しく取得できたか確認
	expectedTodoCount := 2 // ダミーデータが3つ挿入されていると仮定
	if len(myCommits) != expectedTodoCount {
		t.Errorf("Expected %d TODO items, got %d", expectedTodoCount, len(myCommits))
	}

	// ここでtodosを期待するデータと比較して、JSONレスポンスの内容を検証することができます

	// クリーンアップ
	resp.Body.Close()
}

// テスト用のデータベースセットアップ
func setupTestCommitsDatabase() (*gorm.DB, func()) {
	var err error
	db, err = gorm.Open("sqlite3", "test.db") // テスト用のSQLiteデータベースを使用
	if err != nil {
		log.Fatal(err)
	}

	// マイグレーションを実行してテーブルを作成
	db.AutoMigrate(&model.MyCommit{})

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
func createDummyCommitsData() {
	dummyCommits := []model.MyCommit{
		{SHA: "aaaaa", Message: "Message1", Date: time.Now(), Additions: 0, Deletions: 0, Total: 1},
		{SHA: "bbbbb", Message: "Message2", Date: time.Now(), Additions: 0, Deletions: 0, Total: 1},
	}

	for _, commits := range dummyCommits {
		db.Create(&commits)
	}
}
