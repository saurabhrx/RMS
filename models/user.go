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
