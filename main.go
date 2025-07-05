package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Models
type Class struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
}

type Student struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"not null"`
	RollNumber string `json:"roll_number" gorm:"unique;not null"`
	ClassID    uint   `json:"class_id" gorm:"not null"`
	Class      Class  `json:"class" gorm:"foreignKey:ClassID"`
}

var db *gorm.DB

func main() {
	// Database connection
	var err error
	dsn := "postgresql://postgres:fTCNqYhMCcceJvrZJknDNNvxmgRzUXLc@switchyard.proxy.rlwy.net:30216/railway"
	// For Railway, use your DATABASE_URL:
	// dsn := "your_railway_database_url_here"

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&Class{}, &Student{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware for frontend access
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// test running
	r.GET("/riya", hiRiya)

	// Class routes
	r.GET("/classes", getClasses)
	r.GET("/classes/:id", getClass)
	r.POST("/classes", createClass)
	r.PUT("/classes/:id", updateClass)
	r.DELETE("/classes/:id", deleteClass)

	// Student routes
	r.GET("/students", getStudents)
	r.GET("/students/:id", getStudent)
	r.POST("/students", createStudent)
	r.PUT("/students/:id", updateStudent)
	r.DELETE("/students/:id", deleteStudent)

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Class CRUD operations
func getClasses(c *gin.Context) {
	var classes []Class
	db.Find(&classes)
	c.JSON(200, classes)
}

func hiRiya(c *gin.Context) {
	var msg string
	msg = "Hi Riya This API is Working 💜"
	c.JSON(200, msg)
}

func getClass(c *gin.Context) {
	id := c.Param("id")
	var class Class

	if err := db.First(&class, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Class not found"})
		return
	}

	c.JSON(200, class)
}

func createClass(c *gin.Context) {
	var class Class

	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&class).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create class"})
		return
	}

	c.JSON(201, class)
}

func updateClass(c *gin.Context) {
	id := c.Param("id")
	var class Class

	if err := db.First(&class, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Class not found"})
		return
	}

	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Save(&class)
	c.JSON(200, class)
}

func deleteClass(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&Class{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete class"})
		return
	}

	c.JSON(200, gin.H{"message": "Class deleted successfully"})
}

// Student CRUD operations
func getStudents(c *gin.Context) {
	var students []Student
	db.Preload("Class").Find(&students)
	c.JSON(200, students)
}

func getStudent(c *gin.Context) {
	id := c.Param("id")
	var student Student

	if err := db.Preload("Class").First(&student, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(200, student)
}

func createStudent(c *gin.Context) {
	var student Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if class exists
	var class Class
	if err := db.First(&class, student.ClassID).Error; err != nil {
		c.JSON(400, gin.H{"error": "Class not found"})
		return
	}

	if err := db.Create(&student).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create student"})
		return
	}

	// Load the class information
	db.Preload("Class").First(&student, student.ID)
	c.JSON(201, student)
}

func updateStudent(c *gin.Context) {
	id := c.Param("id")
	var student Student

	if err := db.First(&student, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Student not found"})
		return
	}

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if class exists
	var class Class
	if err := db.First(&class, student.ClassID).Error; err != nil {
		c.JSON(400, gin.H{"error": "Class not found"})
		return
	}

	db.Save(&student)
	db.Preload("Class").First(&student, student.ID)
	c.JSON(200, student)
}

func deleteStudent(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&Student{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(200, gin.H{"message": "Student deleted successfully"})
}
