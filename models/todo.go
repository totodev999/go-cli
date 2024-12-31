package models

import "gorm.io/gorm"

type Todo struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	DueDate     string `json:"due_date"`
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Todo{})
}
