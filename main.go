package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func addDummyData(db *sql.DB) error {
	dummyAlbums := []album{
		{ID: "4", Title: "Album Four", Artist: "Artist A", Price: 29.99},
		{ID: "5", Title: "Album Five", Artist: "Artist B", Price: 19.99},
		{ID: "6", Title: "Album Six", Artist: "Artist C", Price: 12.99},
	}

	insertStmt, err := db.Prepare("INSERT INTO albums(id, title, artist, price) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	for _, a := range dummyAlbums {
		_, err := insertStmt.Exec(a.ID, a.Title, a.Artist, a.Price)
		if err != nil {
			return err
		}
	}

	fmt.Println("Dummy data inserted into the database.")
	return nil
}

func main() {
	// データベースに接続
	db, err := sql.Open("sqlite3", "albums.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// albumsテーブルの作成
	// createTable := `
	// 	CREATE TABLE IF NOT EXISTS albums (
	// 		id TEXT PRIMARY KEY,
	// 		title TEXT,
	// 		artist TEXT,
	// 		price REAL
	// 	)
	// `
	// _, err = db.Exec(createTable)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // ダミーデータの追加
	// if err := addDummyData(db); err != nil {
	// 	log.Fatal(err)
	// }

	// Ginルーターのセットアップ
	router := gin.Default()
	router.GET("/albums", getAlbums)
	// 他のエンドポイントのセットアップ

	// サーバーの起動
	router.Run("localhost:8081")
}

func getAlbums(c *gin.Context) {
	// SQLiteデータベースに接続
	db, err := sql.Open("sqlite3", "albums.db")
	if err != nil {
		// エラーハンドリング
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer db.Close()

	// データベースからデータを取得
	rows, err := db.Query("SELECT id, title, artist FROM albums")
	if err != nil {
		// エラーハンドリング
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data from the database"})
		return
	}
	defer rows.Close()

	// 取得したデータをスライスに格納
	var albums []album
	for rows.Next() {
		var a album
		err := rows.Scan(&a.ID, &a.Title, &a.Artist)
		if err != nil {
			// エラーハンドリング
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan data from the database"})
			return
		}
		albums = append(albums, a)
	}

	c.JSON(http.StatusOK, albums)
}
