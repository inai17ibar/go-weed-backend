// db/db.go
package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var instance *gorm.DB

func InitDB() {
	var err error
	instance, err = gorm.Open("sqlite3", "todos.db")
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *gorm.DB {
	return instance
}

func CloseDB() {
	instance.Close()
}
