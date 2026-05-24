package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dratbo/satisfactory-task-manager/gateway/internal/clients"
)

type TaskHandler struct {
	taskClient *clients.TaskClient
	dataClient *clients.DataClient
	tasksTmpl  *template.Template
}

type TaskView struct {
	ID           int64
	Title        string
	Description  string
	Status       string
	CreatedAt    string
	TargetAmount float64
	RecipeName   string
	ProductName  string
	IconURL      string
	Ingredients  []ingredientRow
}

func NewTaskHandler(taskClient *clients.TaskClient, dataClient *clients.DataClient) (*TaskHandler, error) {
	funcMap := template.FuncMap{
		"formatItem": formatItemName,
	}
	tasksTmpl, err := template.New("tasks.html").Funcs(funcMap).ParseFiles("templates/tasks.html")
	if err != nil {
		return nil, err
	}
	return &TaskHandler{
		taskClient: taskClient,
		dataClient: dataClient,
		tasksTmpl:  tasksTmpl,
	}, nil
}

func formatItemName(className string) string {
	s := strings.TrimSuffix(className, "_C")
	s = strings.TrimPrefix(s, "Desc_")
	s = strings.ReplaceAll(s, "_", " ")
	return s
}

func (h *TaskHandler) Index(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("token"); err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	_ = tmpl.Execute(w, nil)
}

func (h *TaskHandler) enrichTask(task clients.Task) TaskView {
	view := TaskView{
		ID:           task.ID,
		Title:        task.Title,
		Description:  task.Description,
		Status:       task.Status,
		CreatedAt:    task.CreatedAt,
		TargetAmount: task.TargetAmount,
	}
	if task.TargetItemClassName == "" {
		return view
	}
	if task.TargetAmount <= 0 {
		view.TargetAmount = 1
	}

	recipe, err := h.dataClient.GetRecipe(task.TargetItemClassName)
	if err != nil || recipe == nil {
		view.RecipeName = formatItemName(task.TargetItemClassName)
		return view
	}

	view.RecipeName = recipe.DisplayName
	for _, ing := range recipe.Ingredients {
		view.Ingredients = append(view.Ingredients, ingredientRow{
			Name:    ing.ItemClassName,
			Amount:  round2(ing.Amount * view.TargetAmount),
			IconURL: clients.ItemIconURL(ing.ItemClassName),
		})
	}
	if len(recipe.Products) > 0 {
		view.ProductName = formatItemName(recipe.Products[0].ItemClassName)
		view.IconURL = clients.ItemIconURL(recipe.Products[0].ItemClassName)
	}
	return view
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tasks, err := h.taskClient.GetTasks(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}

	views := make([]TaskView, 0, len(tasks))
	for _, t := range tasks {
		views = append(views, h.enrichTask(t))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tasksTmpl.Execute(w, views); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	recipeClass := strings.TrimSpace(r.FormValue("target_item_class_name"))
	targetAmount, _ := strconv.ParseFloat(r.FormValue("target_amount"), 64)
	if targetAmount <= 0 {
		targetAmount = 1
	}

	if title == "" && recipeClass != "" {
		if recipe, err := h.dataClient.GetRecipe(recipeClass); err == nil && recipe != nil {
			title = recipe.DisplayName
		}
	}
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	log.Printf("Creating task: title=%s, recipe=%s, amount=%f", title, recipeClass, targetAmount)

	_, err = h.taskClient.CreateTask(cookie.Value, clients.CreateTaskRequest{
		Title:               title,
		Description:         description,
		TargetItemClassName: recipeClass,
		TargetAmount:        targetAmount,
	})
	if err != nil {
		log.Printf("CreateTask error: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	h.GetTasks(w, r)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/tasks/delete/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task id", http.StatusBadRequest)
		return
	}
	if err := h.taskClient.DeleteTask(cookie.Value, id); err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	h.GetTasks(w, r)
}
