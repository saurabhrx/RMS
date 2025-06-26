package dbHelper

import (
	"RMS/database"
	"RMS/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func IsUserExists(email string) (bool, error) {
	query := `SELECT id FROM users WHERE email=$1 AND archived_at IS NULL `
	var id string
	err := database.RMS.Get(&id, query, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil
}

func CreateUser(body *models.UserRequest) (string, error) {
	var (
		userID string
		err    error
	)

	if len(body.Role) == 0 {
		body.Role = []string{"user"}
	}
	if body.CreatedBy != "" {
		query := `INSERT INTO users(name, email, password, created_by) 
		       VALUES ($1, $2, $3, $4) RETURNING id`
		err = database.RMS.QueryRowx(query, body.Name, body.Email, body.Password, body.CreatedBy).Scan(&userID)
		if err != nil {
			return "", err
		}
	} else {
		query := `INSERT INTO users(name, email, password) 
		       VALUES ($1, $2, $3) RETURNING id`
		err = database.RMS.QueryRowx(query, body.Name, body.Email, body.Password).Scan(&userID)
		if err != nil {
			return "", err
		}
	}
	for _, roleType := range body.Role {
		var roleID string
		selectQuery := `SELECT id FROM role WHERE role_type = $1`
		err = database.RMS.Get(&roleID, selectQuery, roleType)
		if err != nil {
			return "", err
		}
		insertQuery := `INSERT INTO user_role(user_id, role_id) VALUES ($1, $2)`
		_, err = database.RMS.Exec(insertQuery, userID, roleID)
		if err != nil {
			return "", err
		}
	}

	for _, address := range body.Address {
		addressQuery := `INSERT INTO user_address(user_id, latitude,longitude) VALUES ($1, $2,$3)`
		_, err = database.RMS.Exec(addressQuery, userID, address.Latitude, address.Longitude)
		if err != nil {
			return "", err
		}
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

func CreateSession(userID string, refreshToken string) (string, error) {
	var sessionID string
	SQL := `INSERT INTO user_session(user_id,refresh_token) VALUES($1,$2) RETURNING id`
	if err := database.RMS.QueryRowx(SQL, userID, refreshToken).Scan(&sessionID); err != nil {
		return "", err
	}
	return sessionID, nil
}

func ValidateSession(userID string, refreshToken string) bool {
	SQL := `SELECT user_id from user_session WHERE refresh_token=$1 AND user_id = $2 AND archived_at IS NULL`
	var user string
	err := database.RMS.Get(&user, SQL, refreshToken, userID)
	if err != nil {
		return false
	}
	return true
}

func Logout(userID string, refreshToken string) error {
	SQL := `UPDATE user_session SET archived_at=now() WHERE user_id=$1 AND refresh_token=$2`
	result, err := database.RMS.Exec(SQL, userID, refreshToken)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no session found to delete")
	}

	return nil
}

func CalculateDistance(userID string, restaurantID string) (models.Distance, error) {
	fmt.Println("Calculating distance")
	var body models.Distance
	query := `SELECT user_address.latitude , user_address.longitude from user_address join users on user_address.user_id = users.id WHERE user_id=$1`
	err := database.RMS.QueryRowx(query, userID).Scan(&body.UserLat, &body.UserLong)
	if err != nil {
		return models.Distance{}, err
	}
	query = `SELECT latitude , longitude from restaurant WHERE id=$1`
	err = database.RMS.QueryRowx(query, restaurantID).Scan(&body.RestaurantLat, &body.RestaurantLong)
	if err != nil {
		return models.Distance{}, err
	}
	return body, nil

}
