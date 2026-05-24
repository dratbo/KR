package models

type Building struct {
	ClassName                string  `json:"class_name"`
	DisplayName              string  `json:"display_name"`
	Description              string  `json:"description"`
	PowerConsumption         float64 `json:"power_consumption"`
	PowerConsumptionExponent float64 `json:"power_consumption_exponent"`
}
