package main

import (
	"log"

	"todo/internal/database"
	"todo/internal/handlers"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	r.GET("/todos", func(c *gin.Context) { handlers.GetTodos(c, db) })
	r.POST("/todos", func(c *gin.Context) { handlers.CreateTodo(c, db) })
	r.GET("/todos/:id", func(c *gin.Context) { handlers.GetTodo(c, db) })
	r.PUT("/todos/:id", func(c *gin.Context) { handlers.UpdateTodo(c, db) })
	r.DELETE("/todos/:id", func(c *gin.Context) { handlers.DeleteTodo(c, db) })

	// Start server
	r.Run(":8080")
}
