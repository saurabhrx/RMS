package models

type UserRequest struct {
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"password" db:"password"`
	Role      []string   `json:"role" db:"role"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	Address   []Location `json:"address" db:"address"`
}

type Location struct {
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}

type LoginRequest struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Distance struct {
	UserLat        float64 `json:"user_lat" db:"user_lat"`
	UserLong       float64 `json:"user_long" db:"user_long"`
	RestaurantLat  float64 `json:"restaurant_lat" db:"restaurant_lat"`
	RestaurantLong float64 `json:"restaurant_long" db:"restaurant_long"`
}

type DistanceRequest struct {
	UserID       string `json:"user_id" db:"user_id"`
	RestaurantID string `json:"restaurant_id" db:"restaurant_id"`
}
