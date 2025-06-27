package handler

import (
	"RMS/database"
	"RMS/database/dbHelper"
	"RMS/middleware"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"github.com/jmoiron/sqlx"
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
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		_, CreateErr := dbHelper.CreateRestaurant(tx, &body)
		return CreateErr

	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create restaurant")
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
	EncodeErr := json.NewEncoder(w).Encode(restaurants)
	if EncodeErr != nil {
		return
	}

}
