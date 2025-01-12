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

type FeedHandlerParams struct {
	dig.In

	AuthHandler *middleware.AuthMiddlewares
	FeedsRepo   *repository.FeedsRepository
	Logger      *structuredlogger.JSONLogger
}
type FeedHandler struct {
	authHandler *middleware.AuthMiddlewares
	feedsRepo   *repository.FeedsRepository
	logger      *structuredlogger.JSONLogger
}

func NewFeedHandler(params FeedHandlerParams) *FeedHandler {
	return &FeedHandler{
		authHandler: params.AuthHandler,
		feedsRepo:   params.FeedsRepo,
		logger:      params.Logger,
	}
}

func (f *FeedHandler) GetFeeds(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	log := &models.LogEntry{
		HTTPRoute: r.URL.Path,
		Timestamp: time.Now(),
	}
	defer f.logger.Log(log, now)

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

	userID, _, ok := f.authHandler.GetUserFromContext(r.Context())
	if !ok {
		err := errors.New("Failed to retrieve user from context")
		log.Error = err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.UserID = userID
	userFeed, err := f.feedsRepo.GetUserFeeds(userID, feedSortingQuery)
	if err != nil {
		log.Error = err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userFeed)
}
