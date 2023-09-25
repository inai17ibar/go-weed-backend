package handler_test

import (
	"go-weed-backend/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCommits(t *testing.T) {
	req, err := http.NewRequest("GET", "/commits?commitCount=2", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	handler.GetCommits(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}
}
