package models

type Schematic struct {
	ClassName      string            `json:"class_name"`
	DisplayName    string            `json:"display_name"`
	Description    string            `json:"description"`
	SchematicType  string            `json:"schematic_type"`
	TimeToComplete float64           `json:"time_to_complete"`
	Costs          []SchematicCost   `json:"costs"`
	Unlocks        []SchematicUnlock `json:"unlocks"`
}

type SchematicCost struct {
	ItemClassName string  `json:"item_class_name"`
	Amount        float64 `json:"amount"`
}

type SchematicUnlock struct {
	Type string `json:"type"` // "Recipe", "Item", "Building", "Schematic"
	Data string `json:"data"` // class name of unlocked entity
}
