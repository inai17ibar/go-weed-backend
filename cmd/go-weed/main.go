package main

import (
	"encoding/json"
	"go-weed-backend/api"
	"go-weed-backend/db"
	"go-weed-backend/router"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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
		api.FetchAndSaveCommits()
		api.FetchAndSaveContribution()
		// タイマーを設定して、一定時間ごとにフェッチ
		time.Sleep(3 * time.Hour) // 例: 6時間ごとにフェッチ
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

	// MongoDBに接続
	dbName := "todo_weed_mongo"
	db.InitDB(dbName) // この接続文字列はconfigから取得 config.MongoDBConnectionString
	defer db.CloseDB()

	// マイグレーションを実行してテーブルを作成
	db.EnsureCollectionExists(db.GetDB().Client(), dbName, "todos")
	db.EnsureCollectionExists(db.GetDB().Client(), dbName, "mycommits")
	db.EnsureCollectionExists(db.GetDB().Client(), dbName, "contributionDays")
	db.EnsureCollectionExists(db.GetDB().Client(), dbName, "taskResults")

	go func() {
		// サーバー起動後、初回のフェッチは遅延させる
		time.Sleep(30 * time.Minute) //もっといい書き方を考えたい、別プログラムとか
		fetchCommitsPeriodically()
	}()

	// ルーターのセットアップ
	r := router.NewRouter()

	// サーバの起動
	log.Printf("Starting server on :%s", config.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.ServerPort, r))
}
