package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserClient struct {
	baseURL string
	client  *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user"`
}

func (c *UserClient) Register(req RegisterRequest) (*AuthResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	log.Printf("Register request to %s: %s", c.baseURL+"/api/register", string(body))

	resp, err := c.client.Post(c.baseURL+"/api/register", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("Register connection error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Register response status: %s", resp.Status)

	if resp.StatusCode != http.StatusCreated {
		// Читаем тело ответа для диагностики
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Printf("Register error body: %s", buf.String())
		return nil, fmt.Errorf("register failed: %s", resp.Status)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}
	return &authResp, nil
}

func (c *UserClient) Login(req LoginRequest) (*AuthResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	log.Printf("Login request to %s: %s", c.baseURL+"/api/login", string(body))

	resp, err := c.client.Post(c.baseURL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("Login connection error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Login response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Printf("Login error body: %s", buf.String())
		return nil, fmt.Errorf("login failed: %s", resp.Status)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}
	return &authResp, nil
}
