package dbHelper

import (
	"RMS/database"
	"RMS/models"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
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

func Register(db sqlx.Ext, body *models.RegisterRequest) (string, error) {
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
		addressQuery := `INSERT INTO user_address(user_id, latitude,longitude) VALUES ($1, $2,$3)`
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

func CreateSession(db sqlx.Ext, userID string, refreshToken string) (string, error) {
	var sessionID string
	SQL := `INSERT INTO user_session(user_id,refresh_token) VALUES($1,$2) RETURNING id`
	if err := db.QueryRowx(SQL, userID, refreshToken).Scan(&sessionID); err != nil {
		return "", err
	}
	return sessionID, nil
}

//func ValidateSession(userID string, refreshToken string) bool {
//	SQL := `SELECT user_id from user_session WHERE refresh_token=$1 AND user_id = $2 AND archived_at IS NULL`
//	var user string
//	err := database.RMS.Get(&user, SQL, refreshToken, userID)
//	if err != nil {
//		return false
//	}
//	return true
//}

func UpdateRefreshToken(db sqlx.Ext, userID, oldToken, newToken string) error {
	query := `UPDATE user_session 
	          SET refresh_token = $1, created_at = NOW()
	          WHERE user_id = $2 AND refresh_token = $3 AND archived_at IS NULL`
	_, err := db.Exec(query, newToken, userID, oldToken)
	return err
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

func GetAllSubadmin(userID string) ([]models.SubadminResponse, error) {
	var subAdmin []models.SubadminResponse
	var roleID string
	SQL := `SELECT id FROM role WHERE role_type='sub-admin'`
	err := database.RMS.Get(&roleID, SQL)
	if err != nil {
		return subAdmin, err
	}
	SQL = `SELECT DISTINCT users.id , users.name FROM users JOIN user_role ON users.id = user_role.user_id 
             WHERE user_role.role_id=$1 AND users.created_by=$2 AND users.archived_at IS NULL`
	err = database.RMS.Select(&subAdmin, SQL, roleID, userID)
	return subAdmin, err
}

func GetUsers(userID string) ([]models.UseResponse, error) {
	var users []models.UseResponse
	SQL := `SELECT DISTINCT users.id , users.name , role.role_type FROM users 
            JOIN user_role ON users.id = user_role.user_id
            JOIN role ON role.id = user_role.role_id where created_by=$1 AND users.archived_at IS NULL `
	err := database.RMS.Select(&users, SQL, userID)
	return users, err

}
