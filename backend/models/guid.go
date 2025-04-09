package models

import (
	"database/sql/driver"

	uuid "github.com/google/uuid"
)

type Guid uuid.UUID

// StringToGuid -> parse string to Guid
func StringToGuid(s string) (Guid, error) {
	id, err := uuid.Parse(s)
	return Guid(id), err
}

// String -> String Representation of Binary16
func (my Guid) String() string {
	return uuid.UUID(my).String()
}

// GormDataType -> sets type to binary(16)
func (my Guid) GormDataType() string {
	return "binary(16)"
}

func (my Guid) MarshalJSON() ([]byte, error) {
	s := uuid.UUID(my)
	str := "\"" + s.String() + "\""
	return []byte(str), nil
}

func (my *Guid) UnmarshalJSON(by []byte) error {
	s, err := uuid.ParseBytes(by)
	*my = Guid(s)
	return err
}

// Scan --> tells GORM how to receive from the database
func (my *Guid) Scan(value interface{}) error {

	bytes, _ := value.([]byte)
	parseByte, err := uuid.FromBytes(bytes)
	*my = Guid(parseByte)
	return err
}

// Value -> tells GORM how to save into the database
func (my Guid) Value() (driver.Value, error) {
	return uuid.UUID(my).MarshalBinary()
}

// func SetId[T *struct{ Id Guid } ](entity T) error {
// func SetId[T ~struct{ Id Guid } ](entity *T) error {
// 	id, err := uuid.NewRandom()
// 	entity.Id = Guid(id)
// 	return err
// }
