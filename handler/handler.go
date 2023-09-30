package handler

import (
	"log"
	"net/http"

	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"go-weed-backend/model"
	"go-weed-backend/util"
)

// dbはデータベースへの参照を保持します。
// 実際のプロジェクトでは、データベースへのアクセス方法をより適切に構造化することが重要です。
var db *gorm.DB

// Initはhandlerパッケージを初期化します。
func Init(database *gorm.DB) {
	db = database
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []model.Todo
	db.Find(&todos)

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := util.ConvertToJapanTime(time.Now()).Truncate(24 * time.Hour)
	todo.Created_Date = currentDate.Format("2006-01-02")

	db.Create(&todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := r.URL.Query().Get("ID")
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing ID parameter")
		return
	}

	var todo model.Todo
	db.First(&todo, todoID)
	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "TODO item not found")
		return
	}

	db.Delete(&todo)
	fmt.Fprintln(w, "TODO item deleted successfully")
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := r.URL.Query().Get("ID")
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing ID parameter")
		return
	}

	var todo model.Todo
	db.First(&todo, todoID)
	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "TODO item not found")
		return
	}

	var updateData struct {
		Title     string `json:"Title"`
		Completed bool   `json:"Completed"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request body")
		return
	}

	todo.Title = updateData.Title
	todo.Completed = updateData.Completed

	db.Save(&todo)
	fmt.Fprintln(w, "TODO item updated successfully")
}

func GetTodosByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("Created_date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing 'Created_date' parameter in the request")
		return
	}

	var todos []model.Todo
	db.Where("created_date = ?", date).Find(&todos)

	// TODOリストをJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func AggregateCommitDataByDate(w http.ResponseWriter, r *http.Request) {
	var commits []model.MyCommit
	if err := db.Find(&commits).Error; err != nil {
		log.Fatalf("Failed to get commits from database: %v", err)
	}

	commitDataByDate := make(map[string]*model.CommitData)
	for _, commit := range commits {
		date := commit.Date.Format("2006-01-02")

		// 日付ごとのデータがなければ新しく作成
		if commitDataByDate[date] == nil {
			commitDataByDate[date] = &model.CommitData{}
		}

		// 日付ごとのコミット数とコード変更量を累計
		commitDataByDate[date].Count++
		commitDataByDate[date].Additions += commit.Additions
		commitDataByDate[date].Deletions += commit.Deletions
		commitDataByDate[date].Total += commit.Total
	}

	// JSONオブジェクトを生成
	jsonData, err := json.Marshal(commitDataByDate)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	// TODOリストをJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonData)
}
