package handler

import (
	"encoding/json"
	"fmt"
	"go-weed-backend/api"
	"net/http"
	"strconv"
)

func GetCommits(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("commitCount")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing 'commitCount' parameter in the request")
		return
	}

	commitCount, err := strconv.Atoi(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "'commitCount' parameter must be a valid integer")
		return
	}

	allCommits, err := api.CallGithubCommitAPI(commitCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error fetching commits")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allCommits)
}
