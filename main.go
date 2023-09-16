package main

import (
	"database/sql"
	"encoding/json"

	"log"
	"net/http"

	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *sql.DB

func main() {
	// データベースに接続
	var err error                             // エラー変数を定義
	db, err = sql.Open("sqlite3", "todos.db") // :=とかくと再宣言することになる
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable := `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			completed BOOLEAN
		)
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// CORSミドルウェアを設定
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),  // Reactアプリケーションのオリジンを指定
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}), // 許可するHTTPメソッドを指定
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			// プリフライトリクエストに対して、必要なCORSヘッダーを返す
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			//w.Header().Set("Access-Control-Allow-Credentials", "true")

			// プリフライトリクエストには空のレスポンスを返す
			w.WriteHeader(http.StatusNoContent)
			return
		}
	})

	http.Handle("/todos", corsHandler(http.HandlerFunc(getTodos)))
	http.Handle("/addTodo", corsHandler(http.HandlerFunc(addTodo)))

	http.ListenAndServe(":8081", nil)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	// データベースからTODOの一覧を取得
	rows, err := db.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
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

	// TODOをデータベースに追加
	_, err := db.Exec("INSERT INTO todos (title, completed) VALUES (?, ?)", todo.Title, todo.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
