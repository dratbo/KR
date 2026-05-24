package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dratbo/satisfactory-task-manager/user-service/internal/config"
	"github.com/dratbo/satisfactory-task-manager/user-service/internal/jwt"
	"github.com/dratbo/satisfactory-task-manager/user-service/internal/models"
	"github.com/dratbo/satisfactory-task-manager/user-service/internal/repository"
)

type AuthHandler struct {
	repo   *repository.UserRepository
	config *config.Config
}

func NewAuthHandler(repo *repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{repo: repo, config: cfg}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "username, email and password required", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.repo.Create(user); err != nil {
		if err == repository.ErrUserAlreadyExists {
			http.Error(w, "username or email already exists", http.StatusConflict)
			return
		}
		log.Println("Create user error:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, h.config.JWTSecret, h.config.JWTExpiresIn)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponse{
		Token: token,
		User:  *user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.FindByUsername(req.Username)
	if err != nil {
		if err == repository.ErrUserNotFound {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !repository.CheckPassword(user.Password, req.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateToken(user.ID, h.config.JWTSecret, h.config.JWTExpiresIn)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponse{
		Token: token,
		User:  *user,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
