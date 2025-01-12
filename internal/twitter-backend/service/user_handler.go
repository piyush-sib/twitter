package service

import (
	"encoding/json"
	"go.uber.org/dig"
	"net/http"

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
}
type UserHandler struct {
	userRepo        *repository.UserRepository
	jwtUtils        *utilities.JWTUtils
	passwordUtils   *utilities.PasswordUtils
	authMiddlewares *middleware.AuthMiddlewares
}

func NewUserHandler(params UserHandlerParams) *UserHandler {
	return &UserHandler{
		userRepo:        params.UserRepo,
		jwtUtils:        params.JwtUtils,
		passwordUtils:   params.PasswordUtils,
		authMiddlewares: params.AuthMiddlewares,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.passwordUtils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	err = h.userRepo.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetUserByEmail(credentials.Email)
	if err != nil || user == nil || !h.passwordUtils.CheckPasswordHash(credentials.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtUtils.GenerateJWT(user.ID, user.Name)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to retrieve user from context", http.StatusInternalServerError)
		return
	}
	_, err := h.userRepo.FollowUserByID(followerData.FollowingUserID, userID)
	if err != nil {
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User followed successfully"})
}
