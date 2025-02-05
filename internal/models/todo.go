package models

import "gorm.io/gorm"

// Todo model with custom ID
type Todo struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'pending'"`
}
