package handler

import (
	"encoding/json"
	"go-weed-backend/internal/model"
	"net/http"
)

func GetContributionDays(w http.ResponseWriter, r *http.Request) {
	var contributionDays []model.ContributionDay
	db.Order("date desc").Find(&contributionDays)

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contributionDays); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
