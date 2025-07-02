package models

type MenuRequest struct {
	Name         string  `json:"name" db:"name"`
	Price        float64 `json:"price" db:"price"`
	CreatedBy    string  `json:"createdBy" db:"created_by"`
	RestaurantId string  `json:"restaurantID" db:"restaurant_id"`
}

type MenuResponse struct {
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
}
