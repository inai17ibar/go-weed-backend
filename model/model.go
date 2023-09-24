package model

import "github.com/jinzhu/gorm"

type Todo struct {
	gorm.Model
	Title        string `json:"Title"`
	Completed    bool   `json:"Completed"`
	Created_Date string `json:"Created_date"`
}
