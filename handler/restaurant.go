package handler

import (
	"RMS/database/dbHelper"
	"RMS/middleware"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"net/http"
)

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var body models.RestaurantRequest
	userID := middleware.UserContext(r)
	body.CreatedBy = userID
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	exists, existsErr := dbHelper.IsRestaurantExists(body.Name, body.Latitude, body.Longitude)
	if existsErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while creating restaurant")
		return
	}
	if exists {
		utils.ResponseError(w, http.StatusConflict, "restaurant already exists")
		return
	}
	_, CreateError := dbHelper.CreateRestaurant(&body)
	if CreateError != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while creating restaurant")
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":    "restaurant successfully created",
		"created_by": userID,
	})
	if EncodeErr != nil {
		return
	}
}

func GetAllRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurants, err := dbHelper.GetAllRestaurants()
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting restaurants")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(restaurants)
	if EncodeErr != nil {
		return
	}
}
