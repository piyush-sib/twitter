package service

import (
	"context"
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
	generateJWT        func(userID int, name string) (string, error)
	hashPassword       func(password string) (string, error)
	checkPasswordHash  func(password, hash string) bool
	getUserFromContext func(ctx context.Context) (userID int, username string, ok bool)
	followUserByID     func(followerID int, userID int) (*models.Followers, error)
	getUserByEmail     func(email string) (*models.User, error)
	createUser         func(user *models.User) error
	logger             func(entry *models.LogEntry, requestTime time.Time)
}

func NewUserHandler(params UserHandlerParams) *UserHandler {
	return &UserHandler{
		generateJWT:        params.JwtUtils.GenerateJWT,
		hashPassword:       params.PasswordUtils.HashPassword,
		checkPasswordHash:  params.PasswordUtils.CheckPasswordHash,
		logger:             params.Logger.Log,
		getUserFromContext: params.AuthMiddlewares.GetUserFromContext,
		createUser:         params.UserRepo.CreateUser,
		getUserByEmail:     params.UserRepo.GetUserByEmail,
		followUserByID:     params.UserRepo.FollowUserByID,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	log := &models.LogEntry{
		HTTPRoute: r.URL.Path,
		Timestamp: time.Now(),
	}
	defer h.logger(log, now)
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

	hashedPassword, err := h.hashPassword(user.Password)
	if err != nil {
		log.Error = err
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	err = h.createUser(&user)
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
	defer h.logger(log, now)
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

	user, err := h.getUserByEmail(credentials.Email)
	if err != nil || user == nil || !h.checkPasswordHash(credentials.Password, user.Password) {
		log.Error = err
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	log.UserID = user.ID

	token, err := h.generateJWT(user.ID, user.Name)
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
	defer h.logger(log, now)
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

	userID, _, ok := h.getUserFromContext(r.Context())
	if !ok {
		err := errors.New("failed to retrieve user from context")
		log.Error = err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.UserID = userID

	_, err := h.followUserByID(followerData.FollowingUserID, userID)
	if err != nil {
		log.Error = err
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User followed successfully"})
}
