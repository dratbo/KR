package repository

import (
	"database/sql"

	"github.com/dratbo/satisfactory-task-manager/satisfactory-data-service/internal/models"
)

type BuildingRepository struct {
	db *sql.DB
}

func NewBuildingRepository(db *sql.DB) *BuildingRepository {
	return &BuildingRepository{db: db}
}

func (r *BuildingRepository) Insert(building *models.Building) error {
	query := `INSERT INTO buildings (class_name, display_name, description, power_consumption, power_consumption_exponent)
	          VALUES ($1, $2, $3, $4, $5) ON CONFLICT (class_name) DO UPDATE SET
	              display_name = EXCLUDED.display_name,
	              description = EXCLUDED.description,
	              power_consumption = EXCLUDED.power_consumption,
	              power_consumption_exponent = EXCLUDED.power_consumption_exponent`
	_, err := r.db.Exec(query, building.ClassName, building.DisplayName, building.Description, building.PowerConsumption, building.PowerConsumptionExponent)
	return err
}

func (r *BuildingRepository) GetAll() ([]models.Building, error) {
	rows, err := r.db.Query(`SELECT class_name, display_name, description, power_consumption, power_consumption_exponent FROM buildings ORDER BY display_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var buildings []models.Building
	for rows.Next() {
		var b models.Building
		err := rows.Scan(&b.ClassName, &b.DisplayName, &b.Description, &b.PowerConsumption, &b.PowerConsumptionExponent)
		if err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}
	return buildings, nil
}
