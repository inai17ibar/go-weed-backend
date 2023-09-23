package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Todo struct {
	gorm.Model
	Title        string `json:"Title"`
	Completed    bool   `json:"Completed"`
	Created_Date string `json:"Created_date"`
}

var db *gorm.DB

func main() {
	// データベースに接続
	var err error
	db, err = gorm.Open("sqlite3", "todos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// マイグレーションを実行してテーブルを作成
	db.AutoMigrate(&Todo{})

	// CORSミドルウェアを設定
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
	})

	http.Handle("/todos", corsHandler(http.HandlerFunc(getTodos)))
	http.Handle("/addTodo", corsHandler(http.HandlerFunc(addTodo)))
	http.Handle("/todos/delete", corsHandler(http.HandlerFunc(deleteTodo)))
	http.Handle("/todos/update", corsHandler(http.HandlerFunc(updateTodo)))
	http.Handle("/todosByDate", corsHandler(http.HandlerFunc(getTodosByDate)))

	http.ListenAndServe(":8081", nil)
}

// データベースから取得したタイムスタンプを日本時間に変換する
func convertToJapanTime(dbTime time.Time) time.Time {
	japanLocation, _ := time.LoadLocation("Asia/Tokyo")
	return dbTime.In(japanLocation)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	db.Find(&todos)

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := convertToJapanTime(time.Now()).Truncate(24 * time.Hour)
	todo.Created_Date = currentDate.Format("2006-01-02")

	db.Create(&todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := r.URL.Query().Get("ID")
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing ID parameter")
		return
	}

	var todo Todo
	db.First(&todo, todoID)
	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "TODO item not found")
		return
	}

	db.Delete(&todo)
	fmt.Fprintln(w, "TODO item deleted successfully")
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := r.URL.Query().Get("ID")
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing ID parameter")
		return
	}

	var todo Todo
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

func getTodosByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("Created_date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing 'Created_date' parameter in the request")
		return
	}

	var todos []Todo
	db.Where("created_date = ?", date).Find(&todos)

	// TODOリストをJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
