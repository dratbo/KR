package main

import (
	"log"
	"net/http"

	"github.com/dratbo/satisfactory-task-manager/gateway/internal/clients"
	"github.com/dratbo/satisfactory-task-manager/gateway/internal/config"
	"github.com/dratbo/satisfactory-task-manager/gateway/internal/handlers"
)

func main() {
	cfg := config.Load()

	userClient := clients.NewUserClient(cfg.UserServiceURL)
	taskClient := clients.NewTaskClient(cfg.TaskServiceURL)
	dataClient := clients.NewDataClient(cfg.DataServiceURL)

	authHandler := handlers.NewAuthHandler(userClient)
	taskHandler, err := handlers.NewTaskHandler(taskClient, dataClient)
	if err != nil {
		log.Fatal("Failed to init task handler:", err)
	}
	recipeHandler, err := handlers.NewRecipeHandler(dataClient)
	if err != nil {
		log.Fatal("Failed to init recipe handler:", err)
	}

	// Статика
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Маршруты аутентификации
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.RegisterPage(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.Register(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.LoginPage(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/logout", authHandler.Logout)

	// Задачи
	http.HandleFunc("/", taskHandler.Index)
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskHandler.GetTasks(w, r)
		} else if r.Method == http.MethodPost {
			taskHandler.CreateTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/tasks/delete/", taskHandler.DeleteTask)
	http.HandleFunc("/recipes/search", recipeHandler.Search)
	http.HandleFunc("/recipes/preview", recipeHandler.Preview)

	log.Printf("Gateway running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
