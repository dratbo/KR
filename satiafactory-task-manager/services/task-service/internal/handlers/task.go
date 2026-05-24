package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dratbo/satisfactory-task-manager/task-service/internal/middleware"
	"github.com/dratbo/satisfactory-task-manager/task-service/internal/models"
	"github.com/dratbo/satisfactory-task-manager/task-service/internal/repository"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateTask handler called")
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		log.Println("Unauthorized: no userID in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("UserID: %d", userID)

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		log.Println("Title is required")
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	task := &models.Task{
		UserID:              userID,
		Title:               req.Title,
		Description:         req.Description,
		Status:              "pending",
		TargetItemClassName: req.TargetItemClassName,
		TargetAmount:        req.TargetAmount,
	}

	log.Printf("Creating task: %+v", task)
	if err := h.repo.Create(task); err != nil {
		log.Printf("Repository create error: %v", err)
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	log.Printf("Task created successfully with ID: %d", task.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.repo.GetByUserID(userID)
	if err != nil {
		http.Error(w, "failed to get tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract task ID from URL path: /tasks/{id}
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteByIDAndUserID(id, userID); err != nil {
		if err == repository.ErrTaskNotFound {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
