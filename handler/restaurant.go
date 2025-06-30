package handler

import (
	"RMS/database/dbHelper"
	"RMS/middleware"
	"RMS/models"
	"RMS/utils"
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
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create restaurant")
		return
	}
	if exists {
		utils.ResponseError(w, http.StatusConflict, "restaurant already exists")
		return
	}
	restaurantID, createErr := dbHelper.CreateRestaurant(&body)

	if createErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create restaurant")
		return
	}

	encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":       "restaurant successfully created",
		"created_by":    userID,
		"restaurant_id": restaurantID,
	})
	if encodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func GetAllRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurants, err := dbHelper.GetAllRestaurants()
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting restaurants")
		return
	}
	encodeErr := json.NewEncoder(w).Encode(restaurants)
	if encodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func GetRestaurantByUerID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	restaurants, err := dbHelper.GetRestaurantByUerID(userID)
	if err != nil && restaurants != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting restaurants")
		return
	}
	if restaurants == nil {
		utils.ResponseError(w, http.StatusOK, "no record available")
		return
	}
	encodeErr := json.NewEncoder(w).Encode(restaurants)
	if encodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}

}
