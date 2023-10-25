// db/db.go
package db

import (
	"fmt"
	"go-weed-backend/internal/model"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var instance *gorm.DB

func InitDB(config model.Config) {
	var err error

	// 環境変数からパスワードを取得
	dbPassword := os.Getenv("DB_MYSQL_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_MYSQL_PASSWORD environment variable is required")
	}

	// データベース名を指定せずにDSNを作成
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, dbPassword, config.DBHost, config.DBPort)

	// backoff用の設定
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute // 最大で2分間再試行

	// backoffを使ってデータベースに接続するロジック
	err = backoff.Retry(func() error {
		var err error
		instance, err = gorm.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return nil
	}, b)

	if err != nil {
		log.Fatalf("Could not connect to database after retries: %v", err)
	}

	// データベースの存在チェックと作成
	instance.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", config.DBName))

	// 作成（または存在する）データベースに再接続
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, dbPassword, config.DBHost, config.DBPort, config.DBName)
	err = backoff.Retry(func() error {
		var err error
		instance, err = gorm.Open("mysql", dsnWithDB)
		if err != nil {
			instance.Close()
			log.Fatal(err)
		}
		return nil
	}, b)

	if err != nil {
		log.Fatalf("Could not connect to database after retries: %v", err)
	}
}

func GetDB() *gorm.DB {
	return instance
}

func CloseDB() {
	instance.Close()
}
