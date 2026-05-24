package handlers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"

	"github.com/dratbo/satisfactory-task-manager/gateway/internal/clients"
)

type RecipeHandler struct {
	dataClient *clients.DataClient
	searchTmpl *template.Template
	previewTmpl *template.Template
}

func NewRecipeHandler(dataClient *clients.DataClient) (*RecipeHandler, error) {
	funcMap := template.FuncMap{"formatItem": formatItemName}
	searchTmpl, err := template.ParseFiles("templates/recipes_search.html")
	if err != nil {
		return nil, err
	}
	previewTmpl, err := template.New("recipe_preview.html").Funcs(funcMap).ParseFiles("templates/recipe_preview.html")
	if err != nil {
		return nil, err
	}
	return &RecipeHandler{
		dataClient:  dataClient,
		searchTmpl:  searchTmpl,
		previewTmpl: previewTmpl,
	}, nil
}

type recipeSearchRow struct {
	ClassName   string
	DisplayName string
	ProductName string
	IconURL     string
}

func (h *RecipeHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) < 2 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<p class="hint">Введите минимум 2 символа для поиска рецепта…</p>`))
		return
	}

	recipes, err := h.dataClient.SearchRecipes(q)
	if err != nil {
		http.Error(w, "Не удалось загрузить рецепты", http.StatusBadGateway)
		return
	}

	rows := make([]recipeSearchRow, 0, len(recipes))
	for _, rec := range recipes {
		row := recipeSearchRow{
			ClassName:   rec.ClassName,
			DisplayName: rec.DisplayName,
		}
		if len(rec.Products) > 0 {
			row.ProductName = rec.Products[0].ItemClassName
			row.IconURL = clients.ItemIconURL(rec.Products[0].ItemClassName)
		}
		rows = append(rows, row)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.searchTmpl.Execute(w, rows)
}

type ingredientRow struct {
	Name    string
	Amount  float64
	IconURL string
}

type recipePreviewData struct {
	RecipeClass string
	Title       string
	IconURL     string
	Multiplier  float64
	Ingredients []ingredientRow
	Products    []ingredientRow
	Duration    float64
}

func (h *RecipeHandler) Preview(w http.ResponseWriter, r *http.Request) {
	className := r.URL.Query().Get("recipe")
	if className == "" {
		http.Error(w, "recipe required", http.StatusBadRequest)
		return
	}

	multiplier, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if multiplier <= 0 {
		multiplier = 1
	}

	recipe, err := h.dataClient.GetRecipe(className)
	if err != nil || recipe == nil {
		http.Error(w, "Рецепт не найден", http.StatusNotFound)
		return
	}

	data := recipePreviewData{
		RecipeClass: recipe.ClassName,
		Title:       recipe.DisplayName,
		Multiplier:  multiplier,
		Duration:    recipe.Duration,
	}

	for _, ing := range recipe.Ingredients {
		data.Ingredients = append(data.Ingredients, ingredientRow{
			Name:    ing.ItemClassName,
			Amount:  round2(ing.Amount * multiplier),
			IconURL: clients.ItemIconURL(ing.ItemClassName),
		})
	}
	for _, prod := range recipe.Products {
		data.Products = append(data.Products, ingredientRow{
			Name:    prod.ItemClassName,
			Amount:  round2(prod.Amount * multiplier),
			IconURL: clients.ItemIconURL(prod.ItemClassName),
		})
		if data.IconURL == "" {
			data.IconURL = clients.ItemIconURL(prod.ItemClassName)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.previewTmpl.Execute(w, data)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
