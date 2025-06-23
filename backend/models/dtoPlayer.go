package models

type DtoPlayer struct {
	Name       string          `json:"name"`
	Attributes []*DtoAttribute `json:"attributes"`
}

func PlayerToDtoPlayer(player *Player) DtoPlayer {
	var attributes []*DtoAttribute
	for _, attr := range player.Attributes {
		if attr.Attribute == "" { // Skip attributes that are bound to another
			attributes = append(attributes, AttributeToDtoAttribute(attr))
		}
	}
	return DtoPlayer{
		Name:       player.Name,
		Attributes: attributes,
	}
}
