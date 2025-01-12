package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"twitter/internal/models"
)

func TestGetTweets_Success(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		getUserTweets: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return []*models.Tweet{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/tweets", nil)
	rr := httptest.NewRecorder()

	handler.GetTweets(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetTweets_InvalidMethod(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		getUserTweets: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return []*models.Tweet{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodPost, "/tweets", nil)
	rr := httptest.NewRecorder()

	handler.GetTweets(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestGetTweets_FailedToRetrieveUser(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 0, "", false
		},
		getUserTweets: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return []*models.Tweet{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/tweets", nil)
	rr := httptest.NewRecorder()

	handler.GetTweets(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestGetTweets_FailedToGetUserTweets(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		getUserTweets: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return nil, errors.New("failed to get tweets")
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/tweets", nil)
	rr := httptest.NewRecorder()

	handler.GetTweets(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestGetTweets_SortingDesc(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		getUserTweets: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return []*models.Tweet{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/tweets?sort=desc", nil)
	rr := httptest.NewRecorder()

	handler.GetTweets(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestPostTweet_Success(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		postTweets: func(userID int, data string) (*models.Tweet, error) {
			return &models.Tweet{Description: data}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	tweet := &models.Tweet{Description: "Hello, world!"}
	body, _ := json.Marshal(tweet)
	req, _ := http.NewRequest(http.MethodPost, "/tweets", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.PostTweet(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestPostTweet_InvalidMethod(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		postTweets: func(userID int, data string) (*models.Tweet, error) {
			return &models.Tweet{Description: data}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/tweets", nil)
	rr := httptest.NewRecorder()

	handler.PostTweet(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestPostTweet_FailedToRetrieveUser(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 0, "", false
		},
		postTweets: func(userID int, data string) (*models.Tweet, error) {
			return &models.Tweet{Description: data}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	tweet := &models.Tweet{Description: "Hello, world!"}
	json.Marshal(tweet)
	req, _ := http.NewRequest(http.MethodPost, "/tweets", http.NoBody)
	rr := httptest.NewRecorder()

	handler.PostTweet(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestPostTweet_InvalidRequestBody(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		postTweets: func(userID int, data string) (*models.Tweet, error) {
			return &models.Tweet{Description: data}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodPost, "/tweets", http.NoBody)
	rr := httptest.NewRecorder()

	handler.PostTweet(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestPostTweet_FailedToPostTweet(t *testing.T) {
	handler := &TweetsHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		postTweets: func(userID int, data string) (*models.Tweet, error) {
			return nil, errors.New("failed to post tweet")
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	tweet := &models.Tweet{Description: "Hello, world!"}
	requestBody, _ := json.Marshal(tweet)
	req, _ := http.NewRequest(http.MethodPost, "/tweets", bytes.NewReader(requestBody))
	rr := httptest.NewRecorder()

	handler.PostTweet(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}
