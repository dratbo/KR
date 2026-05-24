package repository

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/models"
)

func ensureItemStub(tx *sql.Tx, className string) error {
	if className == "" {
		return nil
	}
	display := strings.TrimSuffix(className, "_C")
	display = strings.TrimPrefix(display, "Desc_")
	display = strings.ReplaceAll(display, "_", " ")
	_, err := tx.Exec(
		`INSERT INTO items (class_name, display_name) VALUES ($1, $2) ON CONFLICT (class_name) DO NOTHING`,
		className, display,
	)
	return err
}

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) Insert(recipe *models.Recipe) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	producedInJSON, _ := json.Marshal(recipe.ProducedIn)
	query := `INSERT INTO recipes (class_name, display_name, produced_in, duration, manufactoring_menu_priority)
	          VALUES ($1, $2, $3, $4, $5) ON CONFLICT (class_name) DO UPDATE SET
	              display_name = EXCLUDED.display_name,
	              produced_in = EXCLUDED.produced_in,
	              duration = EXCLUDED.duration,
	              manufactoring_menu_priority = EXCLUDED.manufactoring_menu_priority`
	_, err = tx.Exec(query, recipe.ClassName, recipe.DisplayName, producedInJSON, recipe.Duration, recipe.ManufactoringMenuPriority)
	if err != nil {
		return err
	}

	for _, ing := range recipe.Ingredients {
		if err := ensureItemStub(tx, ing.ItemClassName); err != nil {
			return err
		}
	}
	for _, prod := range recipe.Products {
		if err := ensureItemStub(tx, prod.ItemClassName); err != nil {
			return err
		}
	}

	// Вставка ингредиентов
	for _, ing := range recipe.Ingredients {
		_, err = tx.Exec(`INSERT INTO recipe_ingredients (recipe_class_name, item_class_name, amount) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			recipe.ClassName, ing.ItemClassName, ing.Amount)
		if err != nil {
			return err
		}
	}
	// Вставка продуктов
	for _, prod := range recipe.Products {
		_, err = tx.Exec(`INSERT INTO recipe_products (recipe_class_name, item_class_name, amount) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			recipe.ClassName, prod.ItemClassName, prod.Amount)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *RecipeRepository) loadRecipeDetails(rec *models.Recipe) error {
	ingRows, err := r.db.Query(`SELECT item_class_name, amount FROM recipe_ingredients WHERE recipe_class_name = $1`, rec.ClassName)
	if err != nil {
		return err
	}
	defer ingRows.Close()
	for ingRows.Next() {
		var ing models.Ingredient
		if err := ingRows.Scan(&ing.ItemClassName, &ing.Amount); err != nil {
			return err
		}
		rec.Ingredients = append(rec.Ingredients, ing)
	}

	prodRows, err := r.db.Query(`SELECT item_class_name, amount FROM recipe_products WHERE recipe_class_name = $1`, rec.ClassName)
	if err != nil {
		return err
	}
	defer prodRows.Close()
	for prodRows.Next() {
		var prod models.Product
		if err := prodRows.Scan(&prod.ItemClassName, &prod.Amount); err != nil {
			return err
		}
		rec.Products = append(rec.Products, prod)
	}
	return nil
}

func (r *RecipeRepository) GetByClassName(className string) (*models.Recipe, error) {
	var rec models.Recipe
	var producedInJSON []byte
	err := r.db.QueryRow(
		`SELECT class_name, display_name, produced_in, duration, manufactoring_menu_priority FROM recipes WHERE class_name = $1`,
		className,
	).Scan(&rec.ClassName, &rec.DisplayName, &producedInJSON, &rec.Duration, &rec.ManufactoringMenuPriority)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(producedInJSON, &rec.ProducedIn)
	if err := r.loadRecipeDetails(&rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RecipeRepository) Search(query string, limit int) ([]models.Recipe, error) {
	if limit <= 0 {
		limit = 30
	}
	rows, err := r.db.Query(
		`SELECT class_name, display_name, produced_in, duration, manufactoring_menu_priority
		 FROM recipes
		 WHERE display_name ILIKE '%' || $1 || '%'
		 ORDER BY display_name
		 LIMIT $2`,
		query, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var rec models.Recipe
		var producedInJSON []byte
		if err := rows.Scan(&rec.ClassName, &rec.DisplayName, &producedInJSON, &rec.Duration, &rec.ManufactoringMenuPriority); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(producedInJSON, &rec.ProducedIn)
		if err := r.loadRecipeDetails(&rec); err != nil {
			return nil, err
		}
		recipes = append(recipes, rec)
	}
	return recipes, nil
}

func (r *RecipeRepository) Count() (int, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM recipes`).Scan(&n)
	return n, err
}

func (r *RecipeRepository) GetAll() ([]models.Recipe, error) {
	rows, err := r.db.Query(`SELECT class_name, display_name, produced_in, duration, manufactoring_menu_priority FROM recipes ORDER BY display_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recipes []models.Recipe
	for rows.Next() {
		var rec models.Recipe
		var producedInJSON []byte
		err := rows.Scan(&rec.ClassName, &rec.DisplayName, &producedInJSON, &rec.Duration, &rec.ManufactoringMenuPriority)
		if err != nil {
			return nil, err
		}
		_ = json.Unmarshal(producedInJSON, &rec.ProducedIn)
		if err := r.loadRecipeDetails(&rec); err != nil {
			return nil, err
		}
		recipes = append(recipes, rec)
	}
	return recipes, nil
}
