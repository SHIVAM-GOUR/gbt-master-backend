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

func GetClasses(w http.ResponseWriter, r *http.Request) {
	var classes []models.Class
	config.DB.Find(&classes)

	utils.SendJSONResponse(w, http.StatusOK, classes)
}

func GetClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var class models.Class

	if err := config.DB.First(&class, id).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "Class not found")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, class)
}

func CreateClass(w http.ResponseWriter, r *http.Request) {
	var class models.Class

	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Create(&class).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create class")
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, class)
}

func UpdateClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var class models.Class

	if err := config.DB.First(&class, id).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "Class not found")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse ID from URL param and set it to maintain the same ID
	idInt, _ := strconv.Atoi(id)
	class.ID = uint(idInt)

	config.DB.Save(&class)
	utils.SendJSONResponse(w, http.StatusOK, class)
}

func DeleteClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := config.DB.Delete(&models.Class{}, id).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete class")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Class deleted successfully"})
}
