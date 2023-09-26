package main

import (
	"go-weed-backend/handler"
	"go-weed-backend/model"
	"go-weed-backend/router"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//var db *gorm.DB

func main() {
	// データベースに接続
	//var err error
	db, err := gorm.Open("sqlite3", "todos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// マイグレーションを実行してテーブルを作成
	db.AutoMigrate(&model.Todo{})

	// ハンドラーの初期化
	handler.Init(db)

	// ルーターのセットアップ
	r := router.NewRouter()

	// サーバの起動
	log.Fatal(http.ListenAndServe(":8081", r))

	// commitCount := 5 // 例として、5を用います。
	// commits, err := api.CallGithubCommitAPI(commitCount)
	// if err != nil {
	// 	log.Fatalf("Error calling API: %v", err)
	// }

	// fmt.Printf("Latest %d commits:\n", commitCount)
	// // commitsは []*github.RepositoryCommit 型と仮定
	// for _, commit := range commits {
	// 	if commit == nil {
	// 		fmt.Println("  commit is nil")
	// 		continue
	// 	}

	// 	if commit.SHA == nil || commit.Commit == nil || commit.Commit.Message == nil {
	// 		fmt.Println("  commit.SHA, commit.Commit, or commit.Commit.Message is nil")
	// 		continue
	// 	}

	// 	fmt.Printf("  %s - %s\n", *commit.SHA, *commit.Commit.Message)
	// }
}
