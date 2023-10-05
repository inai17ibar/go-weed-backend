package main

import (
	"encoding/json"
	"go-weed-backend/api"
	"go-weed-backend/internal/handler"
	"go-weed-backend/internal/model"
	"go-weed-backend/router"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

// Config は設定の構造体です。
type Config struct {
	ServerPort string `json:"server_port"`
}

// LoadConfig は設定ファイルを読み込みます。
func LoadConfig(filename string) (Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

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

func main() {
	env := os.Getenv("APP_ENV")
	var config Config
	var err error

	if env == "production" {
		config, err = LoadConfig("config.production.json")
	} else {
		config, err = LoadConfig("config.local.json")
	}

	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// データベースに接続
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
	log.Printf("Starting server on :%s", config.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.ServerPort, r))
}
