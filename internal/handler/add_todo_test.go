package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-weed-backend/internal/model"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddTodo(t *testing.T) {
	//fmt.Println("Starting TestAddTodo")
	// テスト用のデータベースをセットアップ
	client, cleanup := setupTestDatabase()
	defer cleanup()

	// テスト用のHTTPリクエストを作成
	//fmt.Println("Create todo")
	todo := model.Todo{ID: primitive.NewObjectID(), Title: "Task 1", Completed: false}
	body, _ := json.Marshal(todo)
	req := httptest.NewRequest("POST", "/addTodo", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	//fmt.Println("Before calling AddTodo")
	// ハンドラ関数を呼び出し
	AddTodo(w, req)
	//fmt.Println("After calling AddTodo")

	// HTTPレスポンスを取得
	resp := w.Result()

	// レスポンスのステータスコードを確認
	//fmt.Println("Checking response")
	if resp.StatusCode != http.StatusOK {
		t.Logf("Request: %+v", req)
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	//fmt.Println("Fetching data from MongoDB")
	// MongoDBからデータを取得し、期待するデータが存在するか確認
	collection := client.Database("testDB").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todos []model.Todo
	cursor, _ := collection.Find(ctx, bson.M{})
	cursor.All(ctx, &todos)

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo to be added, got %d", len(todos))
	} else if todos[0].Title != "Task 1" || todos[0].Completed {
		t.Errorf("Added todo does not match expected values")
	}
}
