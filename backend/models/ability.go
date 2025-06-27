package models

import "encoding/json"

type Ability struct {
	Name            string  `json:"name"`
	Attribute       string  `json:"attribute"`
	Cost            int     `json:"cost"`
	Chance          float32 `json:"chance"`
	TargetAttribute string  `json:"targetAttribute"`
	TargetDamage    int     `json:"targetDamage"`
	Active          bool    `json:"active"` // Indicates if the ability is currently active
}

func (a *Ability) UnmarshalJSON(text []byte) error {
	type Alias Ability
	activeValue := true
	chanceValue := float32(1)
	targetDamageValue := 1
	aux := Alias{
		Active:       activeValue, // set the default value before parsing JSON
		Chance:       chanceValue,
		TargetDamage: targetDamageValue,
	}
	if err := json.Unmarshal(text, &aux); err != nil {
		return err
	}
	*a = Ability(aux)
	return nil
}
