package handler

import (
	"encoding/json"
	"fmt"
	"go-weed-backend/api"
	"go-weed-backend/internal/model"
	"net/http"

	"github.com/google/uuid"
)

func generateUniqueTaskID() string {
	return uuid.New().String()
}

func FetchAndSaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 新しいタスクの結果のレコードを作成
	taskID := generateUniqueTaskID() // 一意のタスクIDを生成する関数
	task := model.TaskResult{ID: taskID, Status: "in-progress"}
	db.Create(&task)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				// エラーハンドリング
				task.Status = "failure"
				task.Result = fmt.Sprintf("An error occurred: %v", r)
				db.Save(&task)
			}
		}()

		api.FetchAndSaveCommits()
		api.FetchAndSaveContribution()

		task.Status = "success"
		db.Save(&task)
	}()

	w.WriteHeader(http.StatusAccepted) // 202 Accepted: 要求は受け入れられたが、処理はまだ完了していないことを示す
	w.Write([]byte(taskID))            // タスクのIDをクライアントに返す
}

// タスクの結果を確認するためのハンドラ
func CheckTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("taskID")
	if taskID == "" {
		http.Error(w, "taskID is required", http.StatusBadRequest)
		return
	}

	var task model.TaskResult
	db.Where("id = ?", taskID).First(&task)

	if task.Status == "" {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// タスクの結果をJSONとして返す
	json.NewEncoder(w).Encode(task)
}
