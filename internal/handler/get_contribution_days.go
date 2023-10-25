package handler

import (
	"context"
	"encoding/json"
	"go-weed-backend/db"
	"go-weed-backend/internal/model"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func GetContributionDays(w http.ResponseWriter, r *http.Request) {
	collection := db.GetDB().Collection("contributionDays")

	// MongoDBからデータを取得
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Printf("Failed to retrieve contribution days: %v", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var contributionDays []model.ContributionDayDB
	if err = cursor.All(context.TODO(), &contributionDays); err != nil {
		log.Printf("Error while decoding contribution days: %v", err)
		http.Error(w, "Failed to decode data", http.StatusInternalServerError)
		return
	}

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contributionDays); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
