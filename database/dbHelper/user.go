package dbHelper

import (
	"RMS/database"
	"RMS/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func IsUserExists(email string) (bool, error) {
	query := `SELECT count(*)>0 FROM users WHERE email=$1 AND archived_at IS NULL `
	var exists bool
	err := database.RMS.Get(&exists, query, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func RegisterUser(db sqlx.Ext, body *models.RegisterRequest) (string, error) {
	var (
		userID string
		err    error
	)

	query := `INSERT INTO users(name, email, password) 
		       VALUES ($1, $2, $3) RETURNING id`
	err = db.QueryRowx(query, body.Name, body.Email, body.Password).Scan(&userID)
	if err != nil {
		return "", err
	}

	var roleID string
	selectQuery := `SELECT id FROM role WHERE role_type = 'user'`
	err = database.RMS.Get(&roleID, selectQuery)
	if err != nil {
		return "", err
	}
	insertQuery := `INSERT INTO user_role(user_id, role_id) VALUES ($1, $2)`
	_, err = db.Exec(insertQuery, userID, roleID)
	if err != nil {
		return "", err
	}

	for _, address := range body.Address {
		addressQuery := `INSERT INTO user_address(user_id, latitude,longitude) VALUES ($1,$2,$3)`
		_, err = db.Exec(addressQuery, userID, address.Latitude, address.Longitude)
		if err != nil {
			return "", err
		}
	}
	return userID, nil
}

// for admin/subadmin
func CreateUserByAdmin(db sqlx.Ext, body *models.CreateUserRequest) (string, error) {
	var (
		userID string
		err    error
	)
	fmt.Println(body)
	if len(body.Role) == 0 {
		body.Role = []string{"user"}
	}
	query := `INSERT INTO users(name, email, password, created_by) 
		       VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRowx(query, body.Name, body.Email, body.Password, body.CreatedBy).Scan(&userID)
	if err != nil {
		return "", err
	}

	for _, roleType := range body.Role {
		var roleID string
		selectQuery := `SELECT id FROM role WHERE role_type = $1`
		err = database.RMS.Get(&roleID, selectQuery, roleType)
		if err != nil {
			return "", err
		}
		insertQuery := `INSERT INTO user_role(user_id, role_id) VALUES ($1, $2)`
		_, err = db.Exec(insertQuery, userID, roleID)
		if err != nil {
			return "", err
		}
	}
	fmt.Println(userID)
	return userID, nil
}

// subadmin

func CreateUserBySubadmin(db sqlx.Ext, body *models.CreateUserRequestBySubadmin) (string, error) {
	var (
		userID string
		err    error
	)

	query := `INSERT INTO users(name, email, password, created_by) 
		       VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRowx(query, body.Name, body.Email, body.Password, body.CreatedBy).Scan(&userID)
	if err != nil {
		return "", err
	}

	var roleID string
	selectQuery := `SELECT id FROM role WHERE role_type = $1`
	err = database.RMS.Get(&roleID, selectQuery, body.Role)
	if err != nil {
		return "", err
	}
	insertQuery := `INSERT INTO user_role(user_id, role_id) VALUES ($1, $2)`
	_, err = db.Exec(insertQuery, userID, roleID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func ValidateUser(email, password string) (string, error) {
	SQL := `Select id , password from users where archived_at IS NULL and email=$1`
	var userId string
	var hashPassword string
	err := database.RMS.QueryRowx(SQL, email).Scan(&userId, &hashPassword)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	passwordErr := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if passwordErr != nil {
		return "", passwordErr
	}
	return userId, nil
}
func CreateAddress(db sqlx.Ext, body *models.UserAddress) error {
	for _, address := range body.Address {
		SQL := `INSERT INTO user_address(user_id,latitude,longitude) VALUES ($1,$2,$3)`
		_, err := db.Exec(SQL, body.UserID, address.Latitude, address.Longitude)
		if err != nil {
			return err
		}
	}
	return nil
}

func CalculateDistance(addressID string, restaurantID string) (models.Distance, error) {
	var body models.Distance
	SQL := `SELECT user_address.latitude , user_address.longitude FROM user_address 
            WHERE id=$1`
	err := database.RMS.QueryRowx(SQL, addressID).Scan(&body.UserLat, &body.UserLong)
	if err != nil {
		return models.Distance{}, err
	}
	SQL = `SELECT latitude , longitude from restaurant WHERE id=$1`
	err = database.RMS.QueryRowx(SQL, restaurantID).Scan(&body.RestaurantLat, &body.RestaurantLong)
	if err != nil {
		return models.Distance{}, err
	}
	return body, nil

}

func GetAllSubadmin(userID string, limit, page int) ([]models.SubadminResponse, error) {
	var subAdmin = make([]models.SubadminResponse, 0)
	var roleID string
	SQL := `SELECT id FROM role WHERE role_type='subadmin'`
	err := database.RMS.Get(&roleID, SQL)
	if err != nil {
		return subAdmin, err
	}
	SQL = `SELECT DISTINCT users.id , users.name FROM users JOIN user_role ON users.id = user_role.user_id 
             WHERE user_role.role_id=$1 AND users.created_by=$2 AND users.archived_at IS NULL ORDER BY users.name LIMIT $3 OFFSET $4`
	err = database.RMS.Select(&subAdmin, SQL, roleID, userID, limit, page)
	return subAdmin, err
}

func GetRoleByUserID(userID string) ([]string, error) {
	var roleID []string
	var roleType []string
	SQL := `SELECT role_id FROM user_role WHERE user_id=$1`
	err := database.RMS.Select(&roleID, SQL, userID)
	if err != nil {
		return []string{}, err
	}
	SQL = `SELECT role_type FROM role WHERE id=ANY($1)`
	err = database.RMS.Select(&roleType, SQL, pq.Array(roleID))
	if err != nil {
		return []string{}, err
	}

	return roleType, nil
}

func GetUsers(userID string, limit, page int) ([]models.UseResponse, error) {
	var users = make([]models.UseResponse, 0)
	SQL := `SELECT DISTINCT users.id , users.name , role.role_type FROM users 
            JOIN user_role ON users.id = user_role.user_id
            JOIN role ON role.id = user_role.role_id where created_by=$1 AND users.archived_at IS NULL 
            ORDER BY users.name LIMIT $2 OFFSET $3`
	err := database.RMS.Select(&users, SQL, userID, limit, page)
	return users, err

}
