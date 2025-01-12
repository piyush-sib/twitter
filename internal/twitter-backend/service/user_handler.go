package service

import (
	"encoding/json"
	"errors"
	"go.uber.org/dig"
	"net/http"
	"time"
	"twitter/internal/structuredlogger"

	"twitter/internal/models"
	"twitter/internal/twitter-backend/middlewares"
	"twitter/internal/twitter-backend/repository"
	"twitter/internal/twitter-backend/utilities"
)

type UserHandlerParams struct {
	dig.In

	UserRepo        *repository.UserRepository
	JwtUtils        *utilities.JWTUtils
	PasswordUtils   *utilities.PasswordUtils
	AuthMiddlewares *middleware.AuthMiddlewares
	Logger          *structuredlogger.JSONLogger
}
type UserHandler struct {
	userRepo        *repository.UserRepository
	jwtUtils        *utilities.JWTUtils
	passwordUtils   *utilities.PasswordUtils
	authMiddlewares *middleware.AuthMiddlewares
	logger          *structuredlogger.JSONLogger
}

func NewUserHandler(params UserHandlerParams) *UserHandler {
	return &UserHandler{
		userRepo:        params.UserRepo,
		jwtUtils:        params.JwtUtils,
		passwordUtils:   params.PasswordUtils,
		logger:          params.Logger,
		authMiddlewares: params.AuthMiddlewares,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	log := &models.LogEntry{
		HTTPRoute: r.URL.Path,
		Timestamp: time.Now(),
	}
	defer h.logger.Log(log, now)
	if r.Method != http.MethodPost {
		err := errors.New("invalid request method")
		log.Error = err
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error = err
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.passwordUtils.HashPassword(user.Password)
	if err != nil {
		log.Error = err
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	err = h.userRepo.CreateUser(&user)
	if err != nil {
		log.Error = err
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	log.UserID = user.ID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	log := &models.LogEntry{
		HTTPRoute: r.URL.Path,
		Timestamp: time.Now(),
	}
	defer h.logger.Log(log, now)
	if r.Method != http.MethodPost {
		err := errors.New("invalid request method")
		log.Error = err
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Error = err
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetUserByEmail(credentials.Email)
	if err != nil || user == nil || !h.passwordUtils.CheckPasswordHash(credentials.Password, user.Password) {
		log.Error = err
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	log.UserID = user.ID

	token, err := h.jwtUtils.GenerateJWT(user.ID, user.Name)
	if err != nil {
		log.Error = err
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	log := &models.LogEntry{
		HTTPRoute: r.URL.Path,
		Timestamp: time.Now(),
	}
	defer h.logger.Log(log, now)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var followerData struct {
		FollowingUserID int `json:"follower_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&followerData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _, ok := h.authMiddlewares.GetUserFromContext(r.Context())
	if !ok {
		err := errors.New("failed to retrieve user from context")
		log.Error = err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.UserID = userID

	_, err := h.userRepo.FollowUserByID(followerData.FollowingUserID, userID)
	if err != nil {
		log.Error = err
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User followed successfully"})
}
