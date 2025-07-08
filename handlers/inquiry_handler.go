package handlers

import (
	"net/http"

	"github.com/SHIVAM-GOUR/gbt-master-backend/config"
	"github.com/SHIVAM-GOUR/gbt-master-backend/models"
	"github.com/SHIVAM-GOUR/gbt-master-backend/utils"
)

func CreateInquiry(w http.ResponseWriter, r *http.Request) {
	var inquiry models.Inquiry

	err := inquiry.FromJSON(r.Body)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Create(&inquiry).Error; err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create inquiry")
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, inquiry)
}
