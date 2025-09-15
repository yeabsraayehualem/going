package example

import (
	"going/internal/database"
)

type ExampleModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

func init() {
	// Register your models here
	database.RegisterModels(&ExampleModel{})
}
