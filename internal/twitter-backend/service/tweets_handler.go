package service

import (
	"encoding/json"
	"errors"
	"go.uber.org/dig"
	"net/http"
	"time"

	"twitter/internal/models"
	"twitter/internal/structuredlogger"
	"twitter/internal/twitter-backend/middlewares"
	"twitter/internal/twitter-backend/repository"
)

type TweetsHandlerParams struct {
	dig.In

	AuthHandler *middleware.AuthMiddlewares
	Logger      *structuredlogger.JSONLogger
	TweetsRepo  *repository.TweetsRepository
}
type TweetsHandler struct {
	authHandler *middleware.AuthMiddlewares
	logger      *structuredlogger.JSONLogger
	tweetsRepo  *repository.TweetsRepository
}

func NewTweetsHandler(params TweetsHandlerParams) *TweetsHandler {
	return &TweetsHandler{
		authHandler: params.AuthHandler,
		logger:      params.Logger,
		tweetsRepo:  params.TweetsRepo,
	}
}

func (t *TweetsHandler) GetTweets(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	log := &models.LogEntry{
		HTTPRoute: r.URL.Path,
		Timestamp: time.Now(),
	}
	defer t.logger.Log(log, now)

	if r.Method != http.MethodGet {
		err := errors.New("Invalid request method")
		log.Error = err

		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	feedSortingQuery := "asc"
	feedSort := r.URL.Query().Get("sort")
	if feedSort == "desc" {
		feedSortingQuery = "desc"
	}

	userID, _, ok := t.authHandler.GetUserFromContext(r.Context())
	if !ok {
		err := errors.New("Failed to retrieve user from context")
		log.Error = err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.UserID = userID
	userFeed, err := t.tweetsRepo.GetUserTweets(userID, feedSortingQuery)
	if err != nil {
		log.Error = err
		http.Error(w, "Failed to retrieve user tweets", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userFeed)
}

func (t *TweetsHandler) PostTweet(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	log := &models.LogEntry{
		HTTPRoute: r.URL.Path,
		Timestamp: time.Now(),
	}
	defer t.logger.Log(log, now)

	if r.Method != http.MethodPost {
		err := errors.New("Invalid request method")
		log.Error = err

		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	userID, _, ok := t.authHandler.GetUserFromContext(r.Context())
	if !ok {
		err := errors.New("Failed to retrieve user from context")
		log.Error = err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tweet models.Tweet
	if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
		log.Error = err
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.UserID = userID
	finalTweet, err := t.tweetsRepo.PostTweets(userID, tweet.Description)
	if err != nil {
		log.Error = err
		http.Error(w, "Failed to post tweet", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(finalTweet)
}
