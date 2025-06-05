package models

// import "text-adventure/types"

type Attribute map[string]int

type Player struct {
	// Id         types.Guid   `gorm:"primaryKey;" json:"id"`
	Name       string       `json:"name"`
	Items      []*Item      `json:"items"`
	Attributes []*Attribute `json:"attributes"`
	Abilities  []*Attribute `json:"abilities"`
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
