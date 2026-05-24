package models

import "time"

type Task struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	// Новые поля
	TargetItemClassName string  `json:"target_item_class_name,omitempty"`
	TargetAmount        float64 `json:"target_amount,omitempty"`
}

type CreateTaskRequest struct {
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	TargetItemClassName string  `json:"target_item_class_name"`
	TargetAmount        float64 `json:"target_amount"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}
