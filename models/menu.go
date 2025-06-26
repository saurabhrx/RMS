package models

type MenuRequest struct {
	Name         string  `json:"name" db:"name"`
	Price        float64 `json:"price" db:"price"`
	CreatedBy    string  `json:"created_by" db:"created_by"`
	RestaurantId string  `json:"restaurant_id" db:"restaurant_id"`
}

type MenuResponse struct {
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
}
