package service

import (
	"encoding/json"
	"go.uber.org/dig"
	"net/http"

	"twitter/internal/models"
	"twitter/internal/twitter-backend/middlewares"
	"twitter/internal/twitter-backend/repository"
)

type TweetsHandlerParams struct {
	dig.In

	AuthHandler *middleware.AuthMiddlewares
	TweetsRepo  *repository.TweetsRepository
}
type TweetsHandler struct {
	authHandler *middleware.AuthMiddlewares
	tweetsRepo  *repository.TweetsRepository
}

func NewTweetsHandler(params TweetsHandlerParams) *TweetsHandler {
	return &TweetsHandler{
		authHandler: params.AuthHandler,
		tweetsRepo:  params.TweetsRepo,
	}
}

func (t *TweetsHandler) GetTweets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	feedSortingQuery := "asc"
	feedSort := r.URL.Query().Get("sort")
	if feedSort == "desc" {
		feedSortingQuery = "desc"
	}

	userID, _, ok := t.authHandler.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve user from context", http.StatusInternalServerError)
		return
	}
	userFeed, err := t.tweetsRepo.GetUserTweets(userID, feedSortingQuery)
	if err != nil {
		http.Error(w, "Failed to retrieve user tweets", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userFeed)
}

func (t *TweetsHandler) PostTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, _, ok := t.authHandler.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve user from context", http.StatusInternalServerError)
		return
	}
	var tweet models.Tweet
	if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	finalTweet, err := t.tweetsRepo.PostTweets(userID, tweet.Description)
	if err != nil {
		http.Error(w, "Failed to post tweet", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(finalTweet)
}
