package dbHelper

import (
	"RMS/database"
	"RMS/models"
	"database/sql"
	"errors"
)

func IsRestaurantExists(name string, lat, long float64) (bool, error) {
	query := `SELECT id FROM restaurant WHERE name=$1 AND latitude=$2 AND longitude=$3 AND archived_at IS NULL `
	var id string
	err := database.RMS.Get(&id, query, name, lat, long)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil

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
func GetAllRestaurants() ([]models.RestaurantResponse, error) {
	query := `SELECT id , name  , contact , latitude , longitude , opening_time,closing_time,created_by 
              FROM restaurant WHERE archived_at IS NULL`
	var restaurants []models.RestaurantResponse
	err := database.RMS.Select(&restaurants, query)
	return restaurants, err
}
func GetRestaurantByUerID(userID string) ([]models.RestaurantResponse, error) {
	query := `SELECT id , name  , contact , latitude , longitude , opening_time,closing_time,created_by 
              FROM restaurant where created_by=$1 AND archived_at IS NULL`
	var restaurants []models.RestaurantResponse
	err := database.RMS.Select(&restaurants, query, userID)
	return restaurants, err
}
