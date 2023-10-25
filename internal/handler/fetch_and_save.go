package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-weed-backend/api"
	"go-weed-backend/db"
	"go-weed-backend/internal/model"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func generateUniqueTaskID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func FetchAndSaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	collection := db.GetDB().Collection("tasks")

	// 新しいタスクの結果のレコードを作成
	taskID := generateUniqueTaskID()
	task := model.TaskResult{ID: taskID, Status: "in-progress"}

	_, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				// エラーハンドリング
				task.Status = "failure"
				task.Result = fmt.Sprintf("An error occurred: %v", r)
				collection.UpdateOne(context.TODO(), bson.M{"_id": taskID}, bson.M{"$set": task})
			}
		}()

		api.FetchAndSaveCommits()
		api.FetchAndSaveContribution()

		task.Status = "success"
		collection.UpdateOne(context.TODO(), bson.M{"_id": taskID}, bson.M{"$set": task})
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(taskID))
}

func CheckTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("taskID")
	if taskID == "" {
		http.Error(w, "taskID is required", http.StatusBadRequest)
		return
	}

	collection := db.GetDB().Collection("tasks")

	var task model.TaskResult
	err := collection.FindOne(context.TODO(), bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}
