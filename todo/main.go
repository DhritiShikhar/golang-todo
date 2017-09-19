// A simple todo application
// Each todo contains a title and status

package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type todo struct {
	gorm.Model        // Contains ID, CreatedAt, UpdatedAt, DeletedAt fields
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Done       bool   `json:"done"`
}

// Database sets up database connection
func Database() *gorm.DB {
	db, err := gorm.Open("postgres", "sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	// Database Migration
	db := Database()
	db.AutoMigrate(&todo{})

	router := gin.Default()

	api := router.Group("/api/todo")
	{
		api.POST("/", Create) // Create todo
		api.GET("/", List)    // Get all todos
		api.POST("/", Load)   // Get one todo
		api.POST("/", Update) // Update todo
		api.POST("/", Delete) // Delete todo
	}

	router.Run()
}

// Create creates a todo
// gin context contains the request data
func Create(c *gin.Context) {
	todo := todo{Title: c.PostForm("title"), Done: false}
	db := Database()
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Todo created",
		"todoid":  todo.ID,
	})
}

// Load loads a single todo
func Load(c *gin.Context) {
	var todo todo
	id := c.Param("Id")

	db := Database()
	db.First(&todo, id)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo Not Found",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   todo,
	})
}

// List lists all todos
func List(c *gin.Context) {
	var todos []todo

	db := Database()
	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo Not Found",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   todos,
	})
}

// Update updates a single todo
func Update(c *gin.Context) {
	var todo todo
	id := c.Param("Id")

	db := Database()
	db.First(&todo, id)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo Not Found",
		})
	}

	if strings.TrimSpace(c.PostForm("title")) != "" {
		db.Model(&todo).Update("title", c.PostForm("title"))
	}

	if c.GetBool("done") {
		db.Model(&todo).Update("done", c.GetBool("done"))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo Updated",
	})
}

// Delete deletes a single todo
func Delete(c *gin.Context) {
	var todo todo
	id := c.Param("Id")

	db := Database()
	db.First(&todo, id)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo Not Found",
		})
	}

	db.Delete(&todo)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo Deleted",
	})
}
