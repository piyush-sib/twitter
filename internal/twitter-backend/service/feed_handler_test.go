package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"twitter/internal/models"
)

func TestGetFeeds_Success(t *testing.T) {
	handler := &FeedHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		getUserFeeds: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return []*models.Tweet{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/feeds", nil)
	rr := httptest.NewRecorder()

	handler.GetFeeds(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetFeeds_InvalidMethod(t *testing.T) {
	handler := &FeedHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		getUserFeeds: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return []*models.Tweet{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodPost, "/feeds", nil)
	rr := httptest.NewRecorder()

	handler.GetFeeds(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestGetFeeds_FailedToRetrieveUser(t *testing.T) {
	handler := &FeedHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 0, "", false
		},
		getUserFeeds: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return []*models.Tweet{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/feeds", nil)
	rr := httptest.NewRecorder()

	handler.GetFeeds(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestGetFeeds_FailedToGetUserFeeds(t *testing.T) {
	handler := &FeedHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		getUserFeeds: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return nil, errors.New("failed to get feeds")
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/feeds", nil)
	rr := httptest.NewRecorder()

	handler.GetFeeds(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestGetFeeds_SortingDesc(t *testing.T) {
	handler := &FeedHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		getUserFeeds: func(userID int, sortingType string) ([]*models.Tweet, error) {
			return []*models.Tweet{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/feeds?sort=desc", nil)
	rr := httptest.NewRecorder()

	handler.GetFeeds(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
