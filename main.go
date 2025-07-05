package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	// Initialize Chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// special API
	r.Route("/riya", func(r chi.Router) {
		r.Get("/", hiRiya)
	})

	// Routes
	r.Route("/classes", func(r chi.Router) {
		r.Get("/", getClasses)
		r.Post("/", createClass)
		r.Get("/{id}", getClass)
		r.Put("/{id}", updateClass)
		r.Delete("/{id}", deleteClass)
	})

	r.Route("/students", func(r chi.Router) {
		r.Get("/", getStudents)
		r.Post("/", createStudent)
		r.Get("/{id}", getStudent)
		r.Put("/{id}", updateStudent)
		r.Delete("/{id}", deleteStudent)
	})

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func hiRiya(w http.ResponseWriter, r *http.Request) {
	var msg string = "Hi Riya Apni backend API is working.. 💜"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// Class CRUD operations
func getClasses(w http.ResponseWriter, r *http.Request) {
	var classes []Class
	db.Find(&classes)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

func getClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var class Class

	if err := db.First(&class, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Class not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}

func createClass(w http.ResponseWriter, r *http.Request) {
	var class Class

	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := db.Create(&class).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create class"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(class)
}

func updateClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var class Class

	if err := db.First(&class, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Class not found"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Parse ID from URL param and set it to maintain the same ID
	idInt, _ := strconv.Atoi(id)
	class.ID = uint(idInt)

	db.Save(&class)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}

func deleteClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := db.Delete(&Class{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete class"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Class deleted successfully"})
}

// Student CRUD operations
func getStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student
	db.Preload("Class").Find(&students)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var student Student

	if err := db.Preload("Class").First(&student, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Student not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Check if class exists
	var class Class
	if err := db.First(&class, student.ClassID).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Class not found"})
		return
	}

	if err := db.Create(&student).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create student"})
		return
	}

	// Load the class information
	db.Preload("Class").First(&student, student.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var student Student

	if err := db.First(&student, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Student not found"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Parse ID from URL param and set it to maintain the same ID
	idInt, _ := strconv.Atoi(id)
	student.ID = uint(idInt)

	// Check if class exists
	var class Class
	if err := db.First(&class, student.ClassID).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Class not found"})
		return
	}

	db.Save(&student)
	db.Preload("Class").First(&student, student.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := db.Delete(&Student{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete student"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Student deleted successfully"})
}
