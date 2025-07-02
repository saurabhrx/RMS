package handler

import (
	"RMS/database/dbHelper"
	"RMS/middleware"
	"RMS/models"
	"RMS/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
		utils.ResponseError(w, http.StatusInternalServerError, "error while creating dishes")
		return
	}
	if exists {
		utils.ResponseError(w, http.StatusConflict, "dish already exists")
		return
	}
	dishID, createErr := dbHelper.CreateDish(&body)

	if createErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create dish")
		return
	}

	encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":    "dish successfully created",
		"created_by": userID,
		"dish_id":    dishID,
	})
	if encodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func GetDishesByRestaurant(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, limitErr := strconv.Atoi(queryParams.Get("page"))
	if limitErr != nil {
		return
	}
	limit, offsetErr := strconv.Atoi(queryParams.Get("limit"))
	if offsetErr != nil {
		return
	}
	vars := mux.Vars(r)
	restaurantId := vars["restaurant_id"]
	dishes, err := dbHelper.GetDishesByRestaurant(restaurantId, limit, page-1)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting dishes")
		return
	}
	encodeErr := json.NewEncoder(w).Encode(dishes)
	if encodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func GetDishesByUserID(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, limitErr := strconv.Atoi(queryParams.Get("page"))
	if limitErr != nil {
		return
	}
	limit, offsetErr := strconv.Atoi(queryParams.Get("limit"))
	if offsetErr != nil {
		return
	}
	userID := middleware.UserContext(r)
	dishes, err := dbHelper.GetDishesByUserID(userID, limit, page-1)
	if err != nil && dishes != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting dishes")
		return
	}
	if dishes == nil {
		utils.ResponseError(w, http.StatusOK, "no record available")
		return
	}
	encodeErr := json.NewEncoder(w).Encode(dishes)
	if encodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}
