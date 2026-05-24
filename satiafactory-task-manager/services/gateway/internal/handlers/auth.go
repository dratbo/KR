package handlers

import (
	"github.com/dratbo/satisfactory-task-manager/gateway/internal/clients"
	"log"
	"net/http"
)

type AuthHandler struct {
	userClient *clients.UserClient
}

func NewAuthHandler(userClient *clients.UserClient) *AuthHandler {
	return &AuthHandler{userClient: userClient}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}
	resp, err := h.userClient.Login(clients.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/register.html")
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, err := h.userClient.Register(clients.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		// Логируем ошибку в консоль
		log.Println("Registration error:", err)
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}
