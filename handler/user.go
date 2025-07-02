package handler

import (
	"RMS/database"
	"RMS/database/dbHelper"
	"RMS/middleware"
	"RMS/models"
	"RMS/utils"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

var json = utils.JSON
var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body models.RegisterRequest
	var userID string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if body.Name == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide name")
		return
	}
	if body.Password == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide password")
		return
	}
	if body.Email == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide email")
		return
	}
	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while creating User")
		return
	}
	if exists {
		utils.ResponseError(w, http.StatusConflict, "user already exists")
		return
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	body.Password = string(hashedPassword)
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		user, createErr := dbHelper.RegisterUser(tx, &body)
		userID = user
		return createErr

	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create new user")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
		"user_id": userID,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func CreateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	var body models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if body.Name == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide name")
		return
	}
	if body.Password == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide password")
		return
	}
	if body.Email == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide email")
		return
	}
	if len(body.Role) == 0 {
		utils.ResponseError(w, http.StatusBadRequest, "please provide at least one role")
		return
	}

	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while creating User")
		return
	}
	if exists {
		utils.ResponseError(w, http.StatusConflict, "user already exists")
		return
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if hashErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	body.Password = string(hashedPassword)
	body.CreatedBy = userID
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		user, createErr := dbHelper.CreateUserByAdmin(tx, &body)
		userID = user
		return createErr

	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create new user")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
		"user_id": userID,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func CreateUserBySubadmin(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	var body models.CreateUserRequestBySubadmin
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if body.Name == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide name")
		return
	}
	if body.Password == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide password")
		return
	}
	if body.Email == "" {
		utils.ResponseError(w, http.StatusBadRequest, "please provide email")
		return
	}
	if body.Role != "user" {
		utils.ResponseError(w, http.StatusBadRequest, "only authorized to create user")
		return
	}

	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while creating User")
		return
	}
	if exists {
		utils.ResponseError(w, http.StatusConflict, "user already exists")
		return
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if hashErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	body.Password = string(hashedPassword)
	body.CreatedBy = userID

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		user, createErr := dbHelper.CreateUserBySubadmin(tx, &body)
		userID = user
		return createErr

	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create new user")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
		"user_id": userID,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to Decode json")
		return
	}
	userID, validateErr := dbHelper.ValidateUser(body.Email, body.Password)
	if validateErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid credentials")
		return
	}
	if userID == "" {
		utils.ResponseError(w, http.StatusBadRequest, "user not found")
		return
	}
	roleType, roleTypeErr := dbHelper.GetRoleByUserID(userID)
	if roleTypeErr != nil {
		fmt.Println(roleTypeErr)
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting role")
		return
	}
	isUser := false
	for _, role := range roleType {
		if role == "user" {
			isUser = true
		}
	}
	if !isUser {
		utils.ResponseError(w, http.StatusBadRequest, "user not found")
		return
	}
	accessToken, accessErr := middleware.GenerateAccessToken(userID, roleType)
	refreshToken, refreshErr := middleware.GenerateRefreshToken(userID)
	if accessErr != nil || refreshErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate jwt token")
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":       "user logged in successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}

}
func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to Decode json")
		return
	}
	userID, validateErr := dbHelper.ValidateUser(body.Email, body.Password)
	if validateErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid credentials")
		return
	}
	if userID == "" {
		utils.ResponseError(w, http.StatusBadRequest, "user not found")
		return
	}
	roleType, roleTypeErr := dbHelper.GetRoleByUserID(userID)
	if roleTypeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting role")
		return
	}
	isAdmin := false
	for _, role := range roleType {
		if role == "admin" {
			isAdmin = true
		}
	}
	if !isAdmin {
		utils.ResponseError(w, http.StatusBadRequest, "user not found")
		return
	}
	accessToken, accessErr := middleware.GenerateAccessToken(userID, roleType)
	refreshToken, refreshErr := middleware.GenerateRefreshToken(userID)
	if accessErr != nil || refreshErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate jwt token")
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":       "user logged in successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}

}
func LoginSubadmin(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to Decode json")
		return
	}
	userID, validateErr := dbHelper.ValidateUser(body.Email, body.Password)
	if validateErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid credentials")
		return
	}
	if userID == "" {
		utils.ResponseError(w, http.StatusBadRequest, "user not found")
		return
	}
	roleType, roleTypeErr := dbHelper.GetRoleByUserID(userID)
	if roleTypeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting role")
		return
	}
	isSubadmin := false
	for _, role := range roleType {
		if role == "subadmin" {
			isSubadmin = true
		}
	}
	if !isSubadmin {
		utils.ResponseError(w, http.StatusBadRequest, "user not found")
		return
	}
	accessToken, accessErr := middleware.GenerateAccessToken(userID, roleType)
	refreshToken, refreshErr := middleware.GenerateRefreshToken(userID)
	if accessErr != nil || refreshErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate jwt token")
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":       "user logged in successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}

}

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	var body models.UserAddress
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to Decode json")
		return
	}
	body.UserID = userID
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		createErr := dbHelper.CreateAddress(tx, &body)
		return createErr

	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create new user")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message": "address created successfully",
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func CalculateDistance(w http.ResponseWriter, r *http.Request) {
	addressID := mux.Vars(r)["address_id"]
	restaurantID := mux.Vars(r)["restaurant_id"]
	locate, err := dbHelper.CalculateDistance(addressID, restaurantID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to calculate distance")
		return
	}
	distance := utils.HaversineDistance(locate)
	EncodeErr := json.NewEncoder(w).Encode(map[string]float64{
		"distance in km": distance,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func GetAllSubadmin(w http.ResponseWriter, r *http.Request) {
	limit, offset := utils.Pagination(r)
	userID := middleware.UserContext(r)
	subAdmin, err := dbHelper.GetAllSubadmin(userID, limit, offset)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get subadmin")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(subAdmin)
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send respond")
	}

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset := utils.Pagination(r)
	userID := middleware.UserContext(r)
	users, err := dbHelper.GetUsers(userID, limit, offset)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get users")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(users)
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	var body models.RefreshToken
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.UserID == "" || body.Token == "" {
		utils.ResponseError(w, http.StatusUnauthorized, "session expired login again")
		return
	}
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(body.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		utils.ResponseError(w, http.StatusUnauthorized, "session expired login again")
		return

	}
	roleType, roleTypeErr := dbHelper.GetRoleByUserID(body.UserID)
	if roleTypeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting role")
		return
	}
	accessToken, err := middleware.GenerateAccessToken(body.UserID, roleType)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate access token")
		return
	}
	refreshToken, RefreshErr := middleware.GenerateRefreshToken(body.UserID)
	if RefreshErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate refresh token")
		return
	}
	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":       "new access token and refresh token generated successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}
