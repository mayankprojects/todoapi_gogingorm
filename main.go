package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Todo model with custom ID
type Todo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'pending'"`
}

var db *gorm.DB
func main() {
	// Database connection
	dsn := "root:Gmail@12@tcp(127.0.0.1:3306)/todoapp"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&Todo{})

	// Create Gin router
	r := gin.Default()

	// Routes
	r.GET("/todos", getTodos)
	r.POST("/todos", createTodo)
	r.GET("/todos/:id", getTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)

	// Start server
	r.Run(":8080")
}

// retrives the query parameter,
// if not query parameter is not present then set the limit value as 10

func getTodos(c *gin.Context) {
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

	var allTodos []Todo

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
		var filteredTodos []Todo
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

// Create a new todo
func createTodo(c *gin.Context) {
	var todo Todo
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

// Get a specific todo
func getTodo(c *gin.Context) {
	var todo Todo
	result := db.First(&todo, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

// Update a todo
func updateTodo(c *gin.Context) {
	var todo Todo
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

// Delete a todo
func deleteTodo(c *gin.Context) {
	result := db.Delete(&Todo{}, c.Param("id"))
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
