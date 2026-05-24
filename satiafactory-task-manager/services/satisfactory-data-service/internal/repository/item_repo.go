package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/models"
	"github.com/lib/pq"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Insert(item *models.Item) error {
	categoriesJSON, _ := json.Marshal(item.Categories)
	query := `INSERT INTO items (class_name, display_name, description, stack_size, energy_value, form, small_icon, big_icon, categories)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	          ON CONFLICT (class_name) DO UPDATE SET
	              display_name = EXCLUDED.display_name,
	              description = EXCLUDED.description,
	              stack_size = EXCLUDED.stack_size,
	              energy_value = EXCLUDED.energy_value,
	              form = EXCLUDED.form,
	              small_icon = EXCLUDED.small_icon,
	              big_icon = EXCLUDED.big_icon,
	              categories = EXCLUDED.categories`
	_, err := r.db.Exec(query, item.ClassName, item.DisplayName, item.Description, item.StackSize,
		item.EnergyValue, item.Form, item.SmallIcon, item.BigIcon, categoriesJSON)
	return err
}

func (r *ItemRepository) GetAll() ([]models.Item, error) {
	rows, err := r.db.Query(`SELECT class_name, display_name, description, stack_size, energy_value, form, small_icon, big_icon, categories FROM items ORDER BY display_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Item
	for rows.Next() {
		var it models.Item
		var categoriesJSON []byte
		err := rows.Scan(&it.ClassName, &it.DisplayName, &it.Description, &it.StackSize, &it.EnergyValue, &it.Form, &it.SmallIcon, &it.BigIcon, &categoriesJSON)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(categoriesJSON, &it.Categories)
		items = append(items, it)
	}
	return items, nil
}

func (r *ItemRepository) GetDisplayNames(classNames []string) (map[string]string, error) {
	result := make(map[string]string, len(classNames))
	if len(classNames) == 0 {
		return result, nil
	}
	rows, err := r.db.Query(`SELECT class_name, display_name FROM items WHERE class_name = ANY($1)`, pq.Array(classNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var className, displayName string
		if err := rows.Scan(&className, &displayName); err != nil {
			return nil, err
		}
		result[className] = displayName
	}
	return result, nil
}

func (r *ItemRepository) GetByClassName(className string) (*models.Item, error) {
	var it models.Item
	var categoriesJSON []byte
	query := `SELECT class_name, display_name, description, stack_size, energy_value, form, small_icon, big_icon, categories FROM items WHERE class_name = $1`
	err := r.db.QueryRow(query, className).Scan(&it.ClassName, &it.DisplayName, &it.Description, &it.StackSize, &it.EnergyValue, &it.Form, &it.SmallIcon, &it.BigIcon, &categoriesJSON)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(categoriesJSON, &it.Categories)
	return &it, nil
}
