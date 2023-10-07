package handler

import (
	"encoding/json"
	"go-weed-backend/internal/model"
	"net/http"
)

func GetCommits(w http.ResponseWriter, r *http.Request) {
	var my_commits []model.MyCommit
	db.Order("date desc").Find(&my_commits)

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(my_commits); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
