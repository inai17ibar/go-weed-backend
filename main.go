package main

import (
	"go-weed-backend/api"
	"go-weed-backend/handler"
	"go-weed-backend/model"
	"log"

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
	//r := router.NewRouter()

	// サーバの起動
	//log.Fatal(http.ListenAndServe(":8081", r))
	api.CallGithubAPI()
}
