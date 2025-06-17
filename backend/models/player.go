package models

import "fmt"

// import "text-adventure/types"

type Player struct {
	// Id         types.Guid   `gorm:"primaryKey;" json:"id"`
	Name       string       `json:"name"`
	Items      []*Item      `json:"items"`
	Attributes []*Attribute `json:"attributes"`
	Abilities  []*Ability   `json:"abilities"`
}

func (player *Player) GetAttributeValue(attributeName string) (int, error) {
	for i := range player.Attributes {
		if (*player.Attributes[i]).Name == attributeName {
			return (*player.Attributes[i]).Value, nil
		}
	}
	return 0, fmt.Errorf("attribute %s not found", attributeName)
}

// getItem retrieves an item by its code from the player's inventory.
func (player *Player) GetItem(itemCode string) *Item {
	for i := range player.Items {
		if player.Items[i].Code == itemCode {
			return player.Items[i]
		}
	}
	return nil
}

// func (p *Player) BeforeCreate(tx *gorm.DB) error {
// 	id, err := uuid.NewRandom()
// 	p.Id = types.Guid(id)
// 	return err
// }
