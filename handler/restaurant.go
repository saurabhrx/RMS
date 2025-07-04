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
	if body.Name == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide name")
		return
	}
	if body.Contact == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide contact")
		return
	}
	if body.Latitude < -90 && body.Longitude > 90 {
		utils.ResponseError(w, http.StatusBadRequest, "latitude must be between -90 and 90 degree")
		return
	}
	if body.Latitude < -180 && body.Longitude > 180 {
		utils.ResponseError(w, http.StatusBadRequest, "latitude must be between -90 and 90 degree")
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
	limit, offset := utils.Pagination(r)
	restaurants, err := dbHelper.GetAllRestaurants(limit, offset)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting restaurants")
		return
	}
	encodeErr := json.NewEncoder(w).Encode(restaurants)
	if encodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func GetRestaurantByUserID(w http.ResponseWriter, r *http.Request) {
	limit, offset := utils.Pagination(r)
	userID := middleware.UserContext(r)
	restaurants, err := dbHelper.GetRestaurantByUserID(userID, limit, offset)
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
