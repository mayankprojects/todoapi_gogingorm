package handlers

import (
	"net/http"
	"strconv"

	"todo/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTodos retrieves todos with optional limit and status filtering
func GetTodos(c *gin.Context, db *gorm.DB) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	if limit < 1 {
		limit = 10
	}

	var todos []models.Todo

	query := db.Limit(limit)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	result := query.Find(&todos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// CreateTodo creates a new todo
func CreateTodo(c *gin.Context, db *gorm.DB) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodo retrieves a specific todo by ID
func GetTodo(c *gin.Context, db *gorm.DB) {
	var todo models.Todo
	result := db.First(&todo, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

// UpdateTodo updates an existing todo
func UpdateTodo(c *gin.Context, db *gorm.DB) {
	var todo models.Todo
	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Save(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo deletes a todo by ID
func DeleteTodo(c *gin.Context, db *gorm.DB) {
	result := db.Delete(&models.Todo{}, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No todo found with given ID"})
		return
	}

	c.Status(http.StatusNoContent)
}
