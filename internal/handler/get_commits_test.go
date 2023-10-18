package handler

import (
	"encoding/json"
	"go-weed-backend/internal/model"
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
	InitForTest(db) // 追加

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
	expectedCommitsCount := 5 // ダミーデータが3つ挿入されていると仮定
	if len(myCommits) != expectedCommitsCount {
		t.Errorf("Expected %d TODO items, got %d", expectedCommitsCount, len(myCommits))
	}
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
		{SHA: "aaaaa", Message: "Message1", Date: time.Now().Add(-24 * time.Hour), Additions: 0, Deletions: 0, Total: 1},
		{SHA: "bbbbb", Message: "Message2", Date: time.Now().Add(-24 * time.Hour), Additions: 1, Deletions: 1, Total: 2},
		{SHA: "ccccc", Message: "Message3", Date: time.Now().Add(-48 * time.Hour), Additions: 2, Deletions: 2, Total: 4},
		{SHA: "ddddd", Message: "Message4", Date: time.Now().Add(-48 * time.Hour), Additions: 3, Deletions: 3, Total: 6},
		{SHA: "eeeee", Message: "Message5", Date: time.Now().Add(-48 * time.Hour), Additions: 4, Deletions: 4, Total: 8},
	}

	for _, commits := range dummyCommits {
		db.Create(&commits)
	}
}
