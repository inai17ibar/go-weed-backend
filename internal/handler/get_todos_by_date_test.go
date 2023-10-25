package handler

import (
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

func TestGetTodosByDate(t *testing.T) {
	// テスト用のデータベースをセットアップ
	client, cleanup := setupTestDatabase()
	defer cleanup()

	// ダミーデータの作成
	date := "2023-09-01"
	dummyTodo := model.Todo{ID: primitive.NewObjectID(), Title: "Task 1", Completed: false, Created_Date: date}
	collection := client.Database("testDB").Collection("todos")
	_, err := collection.InsertOne(context.TODO(), dummyTodo)
	if err != nil {
		t.Fatalf("Failed to insert dummy todo: %v", err)
	}

	req := httptest.NewRequest("GET", "/todosByDate?Created_Date="+date, nil)
	w := httptest.NewRecorder()

	GetTodosByDate(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	cursor, _ := collection.Find(context.TODO(), bson.M{})
	for cursor.Next(context.TODO()) {
		var todo model.Todo
		cursor.Decode(&todo)
		t.Logf("Database entry: %+v", todo)
	}

	var todos []model.Todo
	json.NewDecoder(resp.Body).Decode(&todos)

	if len(todos) != 1 || todos[0].Title != "Task 1" {
		t.Errorf("Expected 1 todo with title 'Task 1', got %+v", todos)
	}
}
