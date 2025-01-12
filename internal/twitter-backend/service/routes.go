package service

import (
	"github.com/gorilla/mux"
	"net/http"
	middleware "twitter/internal/twitter-backend/middlewares"
)

func RegisterRoutes(router *mux.Router, userHandler *UserHandler, auth *middleware.AuthMiddlewares, feedsHandler *FeedHandler, tweetsHandler *TweetsHandler) {
	// Follow user
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/feeds", feedsHandler.GetFeeds)
	protectedMux.HandleFunc("/tweet", tweetsHandler.PostTweet)
	protectedMux.HandleFunc("/follow", userHandler.FollowUser)
	protectedMux.HandleFunc("/get-tweets", tweetsHandler.GetTweets)

	// Apply middleware
	router.Handle("/feeds", auth.AuthMiddleware()(protectedMux)).Methods("GET")
	router.Handle("/tweet", auth.AuthMiddleware()(protectedMux)).Methods("POST")
	router.Handle("/follow", auth.AuthMiddleware()(protectedMux)).Methods("POST")
	router.Handle("/get-tweets", auth.AuthMiddleware()(protectedMux)).Methods("GET")
}
