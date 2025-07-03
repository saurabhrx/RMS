package dbHelper

import (
	"RMS/database"
	"RMS/models"
)

func CreateDish(body *models.MenuRequest) (string, error) {
	query := `INSERT INTO menu(name, price,restaurant_id,created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	var menuID string
	err := database.RMS.QueryRowx(query, body.Name, body.Price, body.RestaurantId, body.CreatedBy).Scan(&menuID)
	if err != nil {
		return "", err
	}
	return menuID, nil
}

func IsDishExists(name string, restaurantId string) (bool, error) {
	query := `SELECT count(*)>0 FROM menu WHERE name=$1 AND restaurant_id=$2 AND archived_at IS NULL `
	var exists bool
	err := database.RMS.Get(&exists, query, name, restaurantId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetDishesByRestaurant(restaurantId string, limit, page int) ([]models.MenuResponse, error) {
	query := `SELECT menu.name, menu.price FROM menu join restaurant on restaurant.id = menu.restaurant_id 
              WHERE restaurant.id=$1 AND restaurant.archived_at IS NULL ORDER BY menu.name LIMIT $2 OFFSET $3`
	var dishes = make([]models.MenuResponse, 0)
	err := database.RMS.Select(&dishes, query, restaurantId, limit, page)
	return dishes, err
}

func GetDishesByUserID(userID string, limit, page int) ([]models.MenuResponse, error) {
	query := `SELECT menu.name, menu.price  FROM menu WHERE created_by=$1 AND archived_at IS NULL 
              ORDER BY menu.name LIMIT $2 OFFSET $3`
	var dishes = make([]models.MenuResponse, 0)
	err := database.RMS.Select(&dishes, query, userID, limit, page)
	return dishes, err
}
