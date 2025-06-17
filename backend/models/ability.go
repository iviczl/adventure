package models

type Ability struct {
	Name      string `json:"name"`
	Attribute string `json:"attribute"`
	Cost      int    `json:"cost"`
	Active    bool   `json:"active"` // Indicates if the ability is currently active
}
