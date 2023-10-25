package handler

import (
	"context"
	"encoding/json"
	"go-weed-backend/db"
	"go-weed-backend/internal/model"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCommits(w http.ResponseWriter, r *http.Request) {
	collection := db.GetDB().Collection("commits")

	// 日付の降順でソートするオプションを設定
	opts := options.Find().SetSort(bson.M{"date": -1}) // -1は降順を意味する

	cur, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		http.Error(w, "Error fetching commits from database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	var my_commits []model.MyCommit
	for cur.Next(context.TODO()) {
		var commit model.MyCommit
		err := cur.Decode(&commit)
		if err != nil {
			log.Println("Error decoding commit:", err)
			continue
		}
		my_commits = append(my_commits, commit)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, "Cursor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(my_commits); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
