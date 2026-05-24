package repository

import (
	"database/sql"
	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/models"
)

type SchematicRepository struct {
	db *sql.DB
}

func NewSchematicRepository(db *sql.DB) *SchematicRepository {
	return &SchematicRepository{db: db}
}

func (r *SchematicRepository) Insert(schematic *models.Schematic) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO schematics (class_name, display_name, description, schematic_type, time_to_complete)
	          VALUES ($1, $2, $3, $4, $5) ON CONFLICT (class_name) DO UPDATE SET
	              display_name = EXCLUDED.display_name,
	              description = EXCLUDED.description,
	              schematic_type = EXCLUDED.schematic_type,
	              time_to_complete = EXCLUDED.time_to_complete`
	_, err = tx.Exec(query, schematic.ClassName, schematic.DisplayName, schematic.Description, schematic.SchematicType, schematic.TimeToComplete)
	if err != nil {
		return err
	}
	// Costs
	for _, cost := range schematic.Costs {
		_, err = tx.Exec(`INSERT INTO schematic_costs (schematic_class_name, item_class_name, amount) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			schematic.ClassName, cost.ItemClassName, cost.Amount)
		if err != nil {
			return err
		}
	}
	// Unlocks
	for _, unlock := range schematic.Unlocks {
		_, err = tx.Exec(`INSERT INTO schematic_unlocks (schematic_class_name, unlock_type, unlock_data) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			schematic.ClassName, unlock.Type, unlock.Data)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
