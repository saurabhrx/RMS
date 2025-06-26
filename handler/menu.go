package handler

import (
	"RMS/database/dbHelper"
	"RMS/middleware"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var body models.MenuRequest
	userID := middleware.UserContext(r)
	body.CreatedBy = userID
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	exists, existsErr := dbHelper.IsDishExists(body.Name, body.RestaurantId)
	if existsErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while creating restaurant")
		return
	}
	if exists {
		utils.ResponseError(w, http.StatusConflict, "restaurant already exists")
		return
	}
	_, CreateError := dbHelper.CreateDish(&body)
	if CreateError != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while creating restaurant")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":    "dish successfully created",
		"created_by": userID,
	})
	if EncodeErr != nil {
		return
	}
}

func GetDishesByRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restaurantId := vars["restaurant_id"]
	dishes, err := dbHelper.GetDishesByRestaurant(restaurantId)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting restaurant")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(dishes)
	if EncodeErr != nil {
		return
	}
}

func GetDishesByUserID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	dishes, err := dbHelper.GetDishesByUserID(userID)
	if err != nil && dishes != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting restaurant")
		return
	}
	if dishes == nil {
		utils.ResponseError(w, http.StatusOK, "no record available")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(dishes)
	if EncodeErr != nil {
		return
	}
}
