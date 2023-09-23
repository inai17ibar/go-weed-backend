package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedDate string    `json:"created_date"`
}

var db *sql.DB

func main() {
	// データベースに接続
	var err error                                               // エラー変数を定義
	db, err = sql.Open("sqlite3", "todos.db?_loc=Asia%2FTokyo") // :=とかくと再宣言することになる
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable := `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			completed BOOLEAN,
			created_at TIMESTAMP DEFAULT (DATETIME(CURRENT_TIMESTAMP,'localtime')),
			created_date TEXT
		)
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// CORSミドルウェアを設定
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),                   // Reactアプリケーションのオリジンを指定
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}), // 許可するHTTPメソッドを指定
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			// プリフライトリクエストに対して、必要なCORSヘッダーを返す
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			// プリフライトリクエストには空のレスポンスを返す
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
	// データベースからTODOの一覧を取得
	rows, err := db.Query("SELECT id, title, completed, created_at FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	// リクエストヘッダーにCORS設定を追加
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // 許可するオリジンを指定

	// リクエストボディからJSONデータを解析
	var todo Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 現在の日付を取得（時間を含まない）
	currentDate := convertToJapanTime(time.Now()).Truncate(24 * time.Hour)
	date := currentDate.Format("2006-01-02")

	// TODOをデータベースに追加
	_, err := db.Exec("INSERT INTO todos (title, completed, created_date) VALUES (?, ?, ?)", todo.Title, todo.Completed, date)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	// リクエストヘッダーにCORS設定を追加
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // 許可するオリジンを指定

	todoID := r.URL.Query().Get("id")

	// todoIDのバリデーションを行うことをお勧めします
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing todoID parameter")
		return
	}

	// データベースからTODO項目を削除
	_, err := db.Exec("DELETE FROM todos WHERE id=?", todoID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "TODO item deleted successfully")
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	// リクエストヘッダーにCORS設定を追加
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // 許可するオリジンを指定
	todoID := r.URL.Query().Get("id")

	// todoIDのバリデーションを行うことをお勧めします
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing todoID parameter")
		return
	}

	// リクエストボディから新しいタイトルを取得
	var updateData struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request body")
		return
	}

	// データベース内のTODOアイテムのタイトルと完了状態を更新
	_, err = db.Exec("UPDATE todos SET title=?, completed=? WHERE id=?", updateData.Title, updateData.Completed, todoID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "TODO item updated successfully")
}

func getTodosByDate(w http.ResponseWriter, r *http.Request) {
	// リクエストヘッダーにCORS設定を追加
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // 許可するオリジンを指定
	date := r.URL.Query().Get("created_date")

	// 日付のバリデーションを行うことをお勧めします
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing 'created_date' parameter in the request")
		return
	}

	// データベースから指定された日付のTODO項目を取得
	rows, err := db.Query("SELECT id, title, completed FROM todos WHERE created_date=?", date)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	// TODOリストをJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
