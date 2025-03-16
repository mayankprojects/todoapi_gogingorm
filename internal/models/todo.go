package models

// Todo model with custom ID`	`
type Todo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'pending'"`
}