package models

type Item struct {
	ClassName   string   `json:"class_name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	StackSize   int      `json:"stack_size"`
	EnergyValue float64  `json:"energy_value"`
	Form        string   `json:"form"`
	SmallIcon   string   `json:"small_icon"`
	BigIcon     string   `json:"big_icon"`
	Categories  []string `json:"categories"`
}
