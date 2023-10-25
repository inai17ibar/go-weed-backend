package handler

import (
	"context"
	"go-weed-backend/db"
	"go-weed-backend/internal/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupTestDatabase() (*mongo.Client, func()) {
	db.InitDB("mongodb://localhost:27017")

	client := db.GetDB().Client() // 既存のMongoDBクライアントインスタンスを取得

	// クリーンアップ用の関数を返す
	cleanup := func() {
		if err := client.Database("testDB").Drop(context.TODO()); err != nil {
			log.Printf("Failed to drop test database: %v", err)
		}
		// if err := client.Disconnect(context.TODO()); err != nil {
		// 	log.Printf("Failed to disconnect from MongoDB: %v", err)
		// }
	}

	return client, cleanup
}

// ダミーのTODOデータを作成しデータベースに挿入
func createDummyTodoData(client *mongo.Client) {
	dummyTodos := []model.Todo{
		{ID: primitive.NewObjectID(), Title: "Task 1", Completed: false, Created_Date: "2023-09-01"},
		{ID: primitive.NewObjectID(), Title: "Task 2", Completed: true, Created_Date: "2023-09-02"},
		{ID: primitive.NewObjectID(), Title: "Task 3", Completed: false, Created_Date: "2023-09-03"},
	}

	collection := client.Database("testDB").Collection("todos")
	for _, todo := range dummyTodos {
		_, err := collection.InsertOne(context.TODO(), todo)
		if err != nil {
			log.Fatalf("Failed to insert dummy todo: %v", err)
		}
	}
}

func setupTestCommitsDatabase() (*mongo.Client, func()) {
	// MongoDBの初期化
	db.InitDB("mongodb://localhost:27017")

	client := db.GetDB().Client() // 既存のMongoDBクライアントインスタンスを取得

	// テスト用のデータベースへの参照を取得
	db := client.Database("testDB")

	// クリーンアップ用の関数を返す
	cleanup := func() {
		if err := db.Drop(context.TODO()); err != nil {
			log.Printf("Failed to drop test database: %v", err)
		}
		// if err := client.Disconnect(context.TODO()); err != nil {
		// 	log.Printf("Failed to disconnect from MongoDB: %v", err)
		// } //接続は次のテストにも持ち越したい
	}

	return client, cleanup
}

// ダミーのcommitデータを作成しデータベースに挿入
func createDummyCommitsData(client *mongo.Client) {
	dummyCommits := []model.MyCommit{
		{SHA: "aaaaa", Message: "Message1", Date: time.Now().Add(-24 * time.Hour), Additions: 0, Deletions: 0, Total: 1},
		{SHA: "bbbbb", Message: "Message2", Date: time.Now().Add(-24 * time.Hour), Additions: 1, Deletions: 1, Total: 2},
		{SHA: "ccccc", Message: "Message3", Date: time.Now().Add(-48 * time.Hour), Additions: 2, Deletions: 2, Total: 4},
		{SHA: "ddddd", Message: "Message4", Date: time.Now().Add(-48 * time.Hour), Additions: 3, Deletions: 3, Total: 6},
		{SHA: "eeeee", Message: "Message5", Date: time.Now().Add(-48 * time.Hour), Additions: 4, Deletions: 4, Total: 8},
	}

	collection := client.Database("testDB").Collection("commits")
	for _, commit := range dummyCommits {
		_, err := collection.InsertOne(context.TODO(), commit)
		if err != nil {
			log.Fatalf("Failed to insert dummy commit: %v", err)
		}
	}
}
