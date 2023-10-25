package handler

import (
	"encoding/json"
	"go-weed-backend/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCommits(t *testing.T) {
	// テスト用のデータベースをセットアップ
	client, cleanup := setupTestCommitsDatabase()
	defer cleanup()

	// ダミーデータを挿入
	createDummyCommitsData(client)

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
