package models

type Recipe struct {
	ClassName                 string       `json:"class_name"`
	DisplayName               string       `json:"display_name"`
	Ingredients               []Ingredient `json:"ingredients"`
	Products                  []Product    `json:"products"`
	ProducedIn                []string     `json:"produced_in"`
	Duration                  float64      `json:"duration"`
	ManufactoringMenuPriority int          `json:"manufactoring_menu_priority"`
}

type Ingredient struct {
	ItemClassName string  `json:"item_class_name"`
	Amount        float64 `json:"amount"`
}

type Product struct {
	ItemClassName string  `json:"item_class_name"`
	Amount        float64 `json:"amount"`
}
