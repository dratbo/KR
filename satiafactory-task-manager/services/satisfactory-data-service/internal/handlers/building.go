package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/repository"
)

type BuildingHandler struct {
	repo *repository.BuildingRepository
}

func NewBuildingHandler(repo *repository.BuildingRepository) *BuildingHandler {
	return &BuildingHandler{repo: repo}
}

func (h *BuildingHandler) ListBuildings(w http.ResponseWriter, r *http.Request) {
	buildings, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildings)
}
