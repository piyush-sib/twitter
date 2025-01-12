package service

import (
	"encoding/json"
	"go.uber.org/dig"
	"net/http"

	"twitter/internal/twitter-backend/middlewares"
	"twitter/internal/twitter-backend/repository"
)

type FeedHandlerParams struct {
	dig.In

	AuthHandler *middleware.AuthMiddlewares
	FeedsRepo   *repository.FeedsRepository
}
type FeedHandler struct {
	authHandler *middleware.AuthMiddlewares
	feedsRepo   *repository.FeedsRepository
}

func NewFeedHandler(params FeedHandlerParams) *FeedHandler {
	return &FeedHandler{
		authHandler: params.AuthHandler,
		feedsRepo:   params.FeedsRepo,
	}
}

func (f *FeedHandler) GetFeeds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	feedSortingQuery := "asc"
	feedSort := r.URL.Query().Get("sort")
	if feedSort == "desc" {
		feedSortingQuery = "desc"
	}

	userID, _, ok := f.authHandler.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve user from context", http.StatusInternalServerError)
		return
	}
	userFeed, err := f.feedsRepo.GetUserFeeds(userID, feedSortingQuery)
	if err != nil {
		http.Error(w, "Failed to retrieve user feeds", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userFeed)
}
