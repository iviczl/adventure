package models

type DtoAttribute struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// DtoAttribute converts an Attribute to a DtoAttribute for serialization.
func AttributeToDtoAttribute(attribute *Attribute) *DtoAttribute {
	return &(DtoAttribute{
		Name:  attribute.Name,
		Value: attribute.Value,
	})
}
