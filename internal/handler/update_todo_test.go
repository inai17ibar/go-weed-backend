package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-weed-backend/internal/model"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateTodo(t *testing.T) {
	// テスト用のデータベースをセットアップ
	client, cleanup := setupTestDatabase()
	defer cleanup()

	collection := client.Database("testDB").Collection("todos")

	// ダミーデータの作成
	dummyTodo := model.Todo{ID: primitive.NewObjectID(), Title: "Task to Update", Completed: false}
	_, err := collection.InsertOne(context.TODO(), dummyTodo)
	if err != nil {
		t.Fatalf("Failed to insert dummy todo: %v", err)
	}

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
	req := httptest.NewRequest("PUT", "/todos/update?ID="+dummyTodo.ID.Hex(), bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// ハンドラの呼び出し
	UpdateTodo(w, req)

	// レスポンスの取得
	resp := w.Result()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// MongoDBからデータを取得し、更新されているか確認
	var updatedTodo model.Todo
	err = collection.FindOne(context.TODO(), bson.M{"_id": dummyTodo.ID}).Decode(&updatedTodo)
	if err != nil {
		t.Fatalf("Failed to fetch updated todo: %v", err)
	}

	if updatedTodo.Title != updateData.Title || updatedTodo.Completed != updateData.Completed {
		t.Errorf("Updated todo does not match expected values. Got: %+v", updatedTodo)
	}
}
