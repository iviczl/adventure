package models

type DtoPlayer struct {
	Name       string       `json:"name"`
	Attributes []*Attribute `json:"attributes"`
}

func PlayerToDtoPlayer(player *Player) DtoPlayer {
	return DtoPlayer{
		Name:       player.Name,
		Attributes: player.Attributes,
	}
}
