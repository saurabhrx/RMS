package models

type RestaurantRequest struct {
	Name        string  `json:"name" db:"name"`
	Contact     string  `json:"contact" db:"contact"`
	Longitude   float64 `json:"longitude" db:"longitude"`
	Latitude    float64 `json:"latitude" db:"latitude"`
	OpeningTime string  `json:"opening_time" db:"opening_time"`
	ClosingTime string  `json:"closing_time" db:"closing_time"`
	CreatedBy   string  `json:"created_by" db:"created_by"`
}

type RestaurantResponse struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Contact     string  `json:"contact" db:"contact"`
	Latitude    float64 `json:"latitude" db:"latitude"`
	Longitude   float64 `json:"longitude" db:"longitude"`
	OpeningTime *string `json:"opening_time" db:"opening_time"`
	ClosingTime *string `json:"closing_time" db:"closing_time"`
	CreatedBy   string  `json:"created_by" db:"created_by"`
}
