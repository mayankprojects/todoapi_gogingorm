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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	status := c.Query("status")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	var allTodos []models.Todo

	// Calculate offset
	offset := (page - 1) * perPage

	// First get the paginated todos without status filter
	result := db.Offset(offset).Limit(perPage).Find(&allTodos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Then filter by status in memory if status parameter provided
	if status != "" {
		var filteredTodos []models.Todo
		for _, todo := range allTodos {
			if todo.Status == status {
				filteredTodos = append(filteredTodos, todo)
			}
		}
		c.JSON(http.StatusOK, filteredTodos)
		return
	}

	// Return all todos if no status filter
	c.JSON(http.StatusOK, allTodos)
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
