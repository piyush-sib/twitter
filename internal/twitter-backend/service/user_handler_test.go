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

func TestRegister_Success(t *testing.T) {
	handler := &UserHandler{
		hashPassword: func(password string) (string, error) {
			return "hashedPassword", nil
		},
		createUser: func(user *models.User) error {
			user.ID = 1
			return nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	user := &models.User{Password: "password"}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestRegister_InvalidMethod(t *testing.T) {
	handler := &UserHandler{
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/register", nil)
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestRegister_InvalidRequestBody(t *testing.T) {
	handler := &UserHandler{
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodPost, "/register", http.NoBody)
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestRegister_FailedToHashPassword(t *testing.T) {
	handler := &UserHandler{
		hashPassword: func(password string) (string, error) {
			return "", errors.New("failed to hash password")
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	user := &models.User{Password: "password"}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestRegister_FailedToCreateUser(t *testing.T) {
	handler := &UserHandler{
		hashPassword: func(password string) (string, error) {
			return "hashedPassword", nil
		},
		createUser: func(user *models.User) error {
			return errors.New("failed to create user")
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	user := &models.User{Password: "password"}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestLogin_Success(t *testing.T) {
	handler := &UserHandler{
		getUserByEmail: func(email string) (*models.User, error) {
			return &models.User{ID: 1, Password: "hashedPassword"}, nil
		},
		checkPasswordHash: func(password, hash string) bool {
			return true
		},
		generateJWT: func(userID int, name string) (string, error) {
			return "token", nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{Email: "test@example.com", Password: "password"}
	body, _ := json.Marshal(credentials)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLogin_InvalidMethod(t *testing.T) {
	handler := &UserHandler{
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	handler := &UserHandler{
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodPost, "/login", http.NoBody)
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestLogin_InvalidEmailOrPassword(t *testing.T) {
	handler := &UserHandler{
		getUserByEmail: func(email string) (*models.User, error) {
			return nil, errors.New("user not found")
		},
		checkPasswordHash: func(password, hash string) bool {
			return false
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{Email: "test@example.com", Password: "password"}
	body, _ := json.Marshal(credentials)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestLogin_FailedToGenerateToken(t *testing.T) {
	handler := &UserHandler{
		getUserByEmail: func(email string) (*models.User, error) {
			return &models.User{ID: 1, Password: "hashedPassword"}, nil
		},
		checkPasswordHash: func(password, hash string) bool {
			return true
		},
		generateJWT: func(userID int, name string) (string, error) {
			return "", errors.New("failed to generate token")
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{Email: "test@example.com", Password: "password"}
	body, _ := json.Marshal(credentials)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestFollowUser_Success(t *testing.T) {
	handler := &UserHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		followUserByID: func(followerID int, userID int) (*models.Followers, error) {
			return &models.Followers{}, nil
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	followerData := struct {
		FollowingUserID int `json:"follower_user_id"`
	}{FollowingUserID: 2}
	body, _ := json.Marshal(followerData)
	req, _ := http.NewRequest(http.MethodPost, "/follow", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.FollowUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestFollowUser_InvalidMethod(t *testing.T) {
	handler := &UserHandler{
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodGet, "/follow", nil)
	rr := httptest.NewRecorder()

	handler.FollowUser(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestFollowUser_InvalidRequestBody(t *testing.T) {
	handler := &UserHandler{
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	req, _ := http.NewRequest(http.MethodPost, "/follow", http.NoBody)
	rr := httptest.NewRecorder()

	handler.FollowUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestFollowUser_FailedToRetrieveUser(t *testing.T) {
	handler := &UserHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 0, "", false
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	followerData := struct {
		FollowingUserID int `json:"follower_user_id"`
	}{FollowingUserID: 2}
	body, _ := json.Marshal(followerData)
	req, _ := http.NewRequest(http.MethodPost, "/follow", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.FollowUser(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestFollowUser_FailedToFollowUser(t *testing.T) {
	handler := &UserHandler{
		getUserFromContext: func(ctx context.Context) (userID int, username string, ok bool) {
			return 1, "testuser", true
		},
		followUserByID: func(followerID int, userID int) (*models.Followers, error) {
			return nil, errors.New("failed to follow user")
		},
		logger: func(entry *models.LogEntry, requestTime time.Time) {},
	}

	followerData := struct {
		FollowingUserID int `json:"follower_user_id"`
	}{FollowingUserID: 2}
	body, _ := json.Marshal(followerData)
	req, _ := http.NewRequest(http.MethodPost, "/follow", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.FollowUser(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}
