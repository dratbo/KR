package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/repository"
)

type ItemHandler struct {
	repo *repository.ItemRepository
}

func NewItemHandler(repo *repository.ItemRepository) *ItemHandler {
	return &ItemHandler{repo: repo}
}

func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	className := r.PathValue("className")
	if className == "" {
		http.Error(w, "missing className", http.StatusBadRequest)
		return
	}
	item, err := h.repo.GetByClassName(className)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
