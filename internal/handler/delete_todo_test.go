package handler

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-weed-backend/internal/model"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestDeleteTodo(t *testing.T) {
	client, cleanup := setupTestDatabase()
	defer cleanup()

	// ダミーデータの作成
	dummyTodo := model.Todo{ID: primitive.NewObjectID(), Title: "Task to Delete", Completed: false}
	collection := client.Database("testDB").Collection("todos")
	insertResult, err := collection.InsertOne(context.TODO(), dummyTodo)
	if err != nil {
		log.Fatalf("Failed to insert dummy todo: %v", err)
	}
	log.Print("insertResult: ", insertResult)

	req := httptest.NewRequest("DELETE", "/todos/delete?ID="+insertResult.InsertedID.(primitive.ObjectID).Hex(), nil)
	w := httptest.NewRecorder()

	DeleteTodo(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// DBからデータを取得し、データが存在しないことを確認
	var todo model.Todo
	err = collection.FindOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}).Decode(&todo)
	// エラーがない場合、つまりドキュメントが見つかった場合
	if err == nil {
		t.Errorf("Expected todo to be deleted, but was found in db")
		return
	}
	// mongo.ErrNoDocuments以外のエラーが発生した場合
	if err != mongo.ErrNoDocuments {
		t.Errorf("Unexpected error occurred: %v", err)
	}
}
