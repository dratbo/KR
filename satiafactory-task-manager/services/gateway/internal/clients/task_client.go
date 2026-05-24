package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TaskClient struct {
	baseURL string
	client  *http.Client
}

func NewTaskClient(baseURL string) *TaskClient {
	return &TaskClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type Task struct {
	ID                  int64   `json:"id"`
	UserID              int64   `json:"user_id"`
	Title               string  `json:"title"`
	Description         string  `json:"description,omitempty"`
	Status              string  `json:"status"`
	CreatedAt           string  `json:"created_at"`
	TargetItemClassName string  `json:"target_item_class_name,omitempty"`
	TargetAmount        float64 `json:"target_amount,omitempty"`
}

type CreateTaskRequest struct {
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	TargetItemClassName string  `json:"target_item_class_name"`
	TargetAmount        float64 `json:"target_amount"`
}

func (c *TaskClient) GetTasks(token string) ([]Task, error) {
	req, _ := http.NewRequest("GET", c.baseURL+"/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get tasks failed: %s", resp.Status)
	}
	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *TaskClient) CreateTask(token string, req CreateTaskRequest) (*Task, error) {
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.baseURL+"/tasks", bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create task failed: %s", resp.Status)
	}
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (c *TaskClient) DeleteTask(token string, taskID int64) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/tasks/%d", c.baseURL, taskID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete task failed: %s", resp.Status)
	}
	return nil
}
