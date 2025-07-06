package handlers

import (
	"net/http"

	"github.com/SHIVAM-GOUR/gbt-master-backend/utils"
)

func HiRiya(w http.ResponseWriter, r *http.Request) {
	msg := "Hi Riya Apni backend API is working.. 💜"
	utils.SendJSONResponse(w, http.StatusOK, msg)
}
