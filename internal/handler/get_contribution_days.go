package handler

import (
	"encoding/json"
	"go-weed-backend/internal/model"
	"net/http"
)

func GetContributionDays(w http.ResponseWriter, r *http.Request) {
	var contributionDaysDB []model.ContributionDayDB
	db.Order("date desc").Find(&contributionDaysDB)

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contributionDaysDB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
