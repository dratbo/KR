package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dratbo/satisfactory-task-manager/task-service/internal/models"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	query := `INSERT INTO tasks (user_id, title, description, status, created_at, target_item_class_name, target_amount) 
              VALUES ($1, $2, $3, $4, NOW(), $5, $6) RETURNING id, created_at`
	err := r.db.QueryRow(query, task.UserID, task.Title, task.Description, task.Status,
		task.TargetItemClassName, task.TargetAmount).Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (r *TaskRepository) GetByUserID(userID int64) ([]models.Task, error) {
	rows, err := r.db.Query(`SELECT id, user_id, title, description, status, created_at, target_item_class_name, target_amount 
                             FROM tasks WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt,
			&t.TargetItemClassName, &t.TargetAmount); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) DeleteByIDAndUserID(taskID, userID int64) error {
	result, err := r.db.Exec(`DELETE FROM tasks WHERE id = $1 AND user_id = $2`, taskID, userID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}
