package utils

import (
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DbClient *gorm.DB

func CreateDbClient() (*gorm.DB, error) {
	var err error
	DbClient, err = gorm.Open(sqlite.Open(filepath.Join(".", "db", "text-adventure.db")), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return nil, err
	}
	return DbClient, nil
}
