package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dratbo/satisfactory-task-manager/task-service/internal/config"
	"github.com/dratbo/satisfactory-task-manager/task-service/internal/database"
	"github.com/dratbo/satisfactory-task-manager/task-service/internal/handlers"
	"github.com/dratbo/satisfactory-task-manager/task-service/internal/middleware"
	"github.com/dratbo/satisfactory-task-manager/task-service/internal/repository"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	instanceID := os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		instanceID = "unknown"
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("database unreachable:", err)
	}
	database.RunMigrations(db, "migrations")

	taskRepo := repository.NewTaskRepository(db)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	// Middleware для добавления заголовка X-Instance-ID
	instanceMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Instance-ID", instanceID)
			next.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks", taskHandler.GetTasks)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	authMiddleware := middleware.AuthMiddleware(cfg)
	handler := instanceMiddleware(authMiddleware(mux))

	log.Printf("Task service (instance %s) running on port %s", instanceID, cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
