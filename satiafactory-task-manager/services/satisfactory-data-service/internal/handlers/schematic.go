package handlers

import (
	"net/http"

	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/repository"
)

type SchematicHandler struct {
	repo *repository.SchematicRepository
}

func NewSchematicHandler(repo *repository.SchematicRepository) *SchematicHandler {
	return &SchematicHandler{repo: repo}
}

func (h *SchematicHandler) ListSchematics(w http.ResponseWriter, r *http.Request) {
	// Можно реализовать, если нужно
	w.WriteHeader(http.StatusNotImplemented)
}
