package models

import "encoding/json"

// map[string]int
type Attribute struct {
	Name      string `json:"name"`
	Attribute string `json:"attribute"`
	Value     int    `json:"value"`
	MaxValue  int    `json:"maxValue"` // Maximum value for the attribute
}

func (a *Attribute) UnmarshalJSON(text []byte) error {
	type Alias Attribute
	alias := Alias{}
	if err := json.Unmarshal(text, &alias); err != nil {
		return err
	}
	if alias.MaxValue == 0 {
		alias.MaxValue = alias.Value // Set MaxValue to Value if not provided
	}
	*a = Attribute(alias)
	return nil
}
