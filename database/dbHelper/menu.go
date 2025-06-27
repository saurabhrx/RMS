package dbHelper

import (
	"RMS/database"
	"RMS/models"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

func CreateDish(db sqlx.Ext, body *models.MenuRequest) (string, error) {
	query := `INSERT INTO menu(name, price,restaurant_id,created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	var menuID string
	err := db.QueryRowx(query, body.Name, body.Price, body.RestaurantId, body.CreatedBy).Scan(&menuID)
	if err != nil {
		return "", err
	}
	return menuID, nil
}

func IsDishExists(name string, restaurantId string) (bool, error) {
	query := `SELECT id FROM menu WHERE name=$1 AND restaurant_id=$2 AND archived_at IS NULL`
	var id string
	err := database.RMS.Get(&id, query, name, restaurantId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil

}

func GetDishesByRestaurant(restaurantId string) ([]models.MenuResponse, error) {
	query := `SELECT menu.name, menu.price FROM menu join restaurant on restaurant.id = menu.restaurant_id 
              WHERE restaurant.id=$1 AND restaurant.archived_at IS NULL `
	var dishes []models.MenuResponse
	err := database.RMS.Select(&dishes, query, restaurantId)
	return dishes, err

}
func GetDishesByUserID(userID string) ([]models.MenuResponse, error) {
	query := `SELECT menu.name, menu.price  FROM menu WHERE created_by=$1 AND archived_at IS NULL`
	var dishes []models.MenuResponse
	err := database.RMS.Select(&dishes, query, userID)
	return dishes, err

}
