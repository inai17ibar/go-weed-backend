package handler

import (
	"encoding/json"
	"go-weed-backend/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAggregateCommitDataByDate(t *testing.T) {

	// テスト用のデータベースをセットアップ
	var cleanup func()
	db, cleanup = setupTestCommitsDatabase()
	defer cleanup()

	// ハンドラの初期化
	Init(db) // 追加

	// テスト用のダミーTODOデータを挿入
	createDummyCommitsData()

	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest("GET", "/commitDataByDate", nil)
	w := httptest.NewRecorder()

	// ハンドラ関数を呼び出し
	AggregateCommitDataByDate(w, req)

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
	var commitDataList []model.CommitData
	err := json.NewDecoder(resp.Body).Decode(&commitDataList)
	if err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}

	// ダミーTODOデータが正しく取得できたか確認
	expectedCommitsCount := 2 // ダミーデータが3つ挿入されていると仮定
	if len(commitDataList) != expectedCommitsCount {
		t.Errorf("Expected %d TODO items, got %d", expectedCommitsCount, len(commitDataList))
	}
	// クリーンアップ
	resp.Body.Close()
}

// // テスト用のデータベースセットアップ
// func setupTestCommitDataDatabase() (*gorm.DB, func()) {
// 	var err error
// 	db, err = gorm.Open("sqlite3", "test.db") // テスト用のSQLiteデータベースを使用
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// マイグレーションを実行してテーブルを作成
// 	db.AutoMigrate(&model.CommitData{})

// 	// クリーンアップ用の関数を返す
// 	cleanup := func() {
// 		db.Close()
// 		if err := os.Remove("test.db"); err != nil {
// 			log.Printf("Failed to remove test database file: %v", err)
// 		}
// 	}

// 	return db, cleanup
// }

// // ダミーのTODOデータを作成しデータベースに挿入
// func createDummyCommitDataByDate() {
// 	dummyData := map[string]model.CommitData{
// 		"2023-10-01": {
// 			Count:     5,
// 			Additions: 50,
// 			Deletions: 20,
// 			Total:     70,
// 		},
// 		"2023-10-02": {
// 			Count:     3,
// 			Additions: 30,
// 			Deletions: 10,
// 			Total:     40,
// 		},
// 		"2023-10-03": {
// 			Count:     2,
// 			Additions: 20,
// 			Deletions: 5,
// 			Total:     25,
// 		},
// 	}

// 	for _, commitData := range dummyData {
// 		db.Create(&commitData)
// 	}
// }
