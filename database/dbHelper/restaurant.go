package dbHelper

import (
	"RMS/database"
	"RMS/models"
)

func IsRestaurantExists(name string, lat, long float64) (bool, error) {
	query := `SELECT count(*)>0 FROM restaurant WHERE name=$1 AND latitude=$2 AND longitude=$3 AND archived_at IS NULL `
	var exists bool
	err := database.RMS.Get(&exists, query, name, lat, long)
	if err != nil {
		return false, err
	}
	return exists, nil

}

func CreateRestaurant(body *models.RestaurantRequest) (string, error) {
	query := `INSERT INTO restaurant(name ,contact,opening_time,closing_time,latitude,longitude,created_by) 
              VALUES ($1, $2, $3, $4,$5,$6,$7) RETURNING id`
	var restID string
	err := database.RMS.QueryRowx(query, body.Name, body.Contact, body.OpeningTime, body.ClosingTime, body.Latitude, body.Longitude, body.CreatedBy).Scan(&restID)
	if err != nil {
		return "", err
	}
	return restID, nil

}
func GetAllRestaurants(limit, page int) ([]models.RestaurantResponse, error) {
	query := `SELECT id , name  , contact , latitude , longitude , opening_time,closing_time,created_by 
              FROM restaurant WHERE archived_at IS NULL ORDER BY name LIMIT $1 OFFSET $2`
	var restaurants = make([]models.RestaurantResponse, 0)
	err := database.RMS.Select(&restaurants, query, limit, page)
	return restaurants, err
}
func GetRestaurantByUserID(userID string, limit, page int) ([]models.RestaurantResponse, error) {
	query := `SELECT id , name  , contact , latitude , longitude , opening_time,closing_time,created_by 
              FROM restaurant where created_by=$1 AND archived_at IS NULL ORDER BY name LIMIT $1 OFFSET $2`
	var restaurants = make([]models.RestaurantResponse, 0)
	err := database.RMS.Select(&restaurants, query, userID, limit, page)
	return restaurants, err
}
