package main

import (
	"go-weed-backend/api"
	"go-weed-backend/handler"
	"go-weed-backend/model"
	"go-weed-backend/router"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func fetchCommitsPeriodically() {
	for {
		fetchAndSaveCommits()
		// タイマーを設定して、一定時間ごとにフェッチ
		time.Sleep(6 * time.Hour) // 例: 6時間ごとにフェッチ
	}
}

func fetchAndSaveCommits() {
	// 最後のコミットの日付を取得
	var lastCommit model.MyCommit
	db.Order("date desc").First(&lastCommit)

	var commits []model.MyCommit

	var err error
	commits, err = api.CallGithubAllCommitAPI()
	if err != nil {
		log.Fatalf("Error fetching commits from GitHub API: %v", err)
	}

	for _, c := range commits {
		// 最後にフェッチしたコミット以降のコミットだけを保存
		if c.Date.After(lastCommit.Date) {
			db.Save(&c)
		}
	}
}

// func getCommits(c *gin.Context) {
// 	var commits []model.MyCommit
// 	db.Order("date desc").Find(&commits)
// 	c.JSON(http.StatusOK, commits)
// }

func main() {
	// データベースに接続
	var err error
	db, err = gorm.Open("sqlite3", "todos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// マイグレーションを実行してテーブルを作成
	db.AutoMigrate(&model.Todo{})
	db.AutoMigrate(&model.MyCommit{}) //これがそのままテーブル名になる

	go func() {
		// サーバー起動後、初回のフェッチは遅延させる
		time.Sleep(60 * time.Minute) //もっといい書き方を考えたい、別プログラムとか
		fetchCommitsPeriodically()
	}()

	// ハンドラーの初期化
	handler.Init(db)

	// ルーターのセットアップ
	r := router.NewRouter()

	// サーバの起動
	log.Fatal(http.ListenAndServe(":8081", r))

	//これをもとにテストをかく
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
