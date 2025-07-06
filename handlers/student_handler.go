package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SHIVAM-GOUR/gbt-master-backend/config"
	"github.com/SHIVAM-GOUR/gbt-master-backend/models"
	"github.com/SHIVAM-GOUR/gbt-master-backend/utils"

	"github.com/go-chi/chi/v5"
)

func GetStudents(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	config.DB.Preload("Class").Find(&students)

	utils.SendJSONResponse(w, http.StatusOK, students)
}

func GetStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var student models.Student

	if err := config.DB.Preload("Class").First(&student, id).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "Student not found")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, student)
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if class exists
	var class models.Class
	if err := config.DB.First(&class, student.ClassID).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Class not found")
		return
	}

	if err := config.DB.Create(&student).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create student")
		return
	}

	// Load the class information
	config.DB.Preload("Class").First(&student, student.ID)

	utils.SendJSONResponse(w, http.StatusCreated, student)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var student models.Student

	if err := config.DB.First(&student, id).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "Student not found")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse ID from URL param and set it to maintain the same ID
	idInt, _ := strconv.Atoi(id)
	student.ID = uint(idInt)

	// Check if class exists
	var class models.Class
	if err := config.DB.First(&class, student.ClassID).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Class not found")
		return
	}

	config.DB.Save(&student)
	config.DB.Preload("Class").First(&student, student.ID)

	utils.SendJSONResponse(w, http.StatusOK, student)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := config.DB.Delete(&models.Student{}, id).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete student")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Student deleted successfully"})
}
