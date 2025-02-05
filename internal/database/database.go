package database

import (
	"todo/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes and returns a database connection
func InitDB() (*gorm.DB, error) {
	// Database connection
	dsn := "root:Gmail@12@tcp(127.0.0.1:3306)/todoapp"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
