package models

const (
	RoleAdmin    = "admin"
	RoleSubadmin = "sub-admin"
	RoleUser     = "user"
)

type RegisterRequest struct {
	Name     string     `json:"name" db:"name"`
	Email    string     `json:"email" db:"email"`
	Password string     `json:"password" db:"password"`
	Address  []Location `json:"address" db:"address"`
}
type CreateUserRequest struct {
	Name      string   `json:"name" db:"name"`
	Email     string   `json:"email" db:"email"`
	Password  string   `json:"password" db:"password"`
	Role      []string `json:"role" db:"role"`
	CreatedBy string   `json:"createdBy" db:"created_by"`
}
type CreateUserRequestBySubadmin struct {
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	Role      string `json:"role" db:"role"`
	CreatedBy string `json:"createdBy" db:"created_by"`
}

type UserAddress struct {
	UserID  string     `json:"userID" db:"user_id"`
	Address []Location `json:"address" db:"address"`
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
	UserLat        float64 `json:"userLatitude" db:"user_latitude"`
	UserLong       float64 `json:"userLongitude" db:"user_longitude"`
	RestaurantLat  float64 `json:"restaurantLatitude" db:"restaurant_latitude"`
	RestaurantLong float64 `json:"restaurantLongitude" db:"restaurant_longitude"`
}

type DistanceRequest struct {
	UserID       string `json:"userID" db:"user_id"`
	RestaurantID string `json:"restaurantID" db:"restaurant_id"`
}

type SubadminResponse struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
type UseResponse struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Role string `json:"role" db:"role_type"`
}

type RefreshToken struct {
	UserID string `json:"userID" db:"user_id"`
	Token  string `json:"refreshToken" db:"refresh_token"`
}
