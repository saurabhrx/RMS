package handler

import (
	"RMS/database/dbHelper"
	"RMS/middleware"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var body models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
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
	fmt.Println(body.Password)
	userID, createErr := dbHelper.CreateUser(&body)
	if createErr != nil || userID == "" {
		fmt.Println(createErr)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create new user")
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
		"user_id": userID,
	})
	if err != nil {
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
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
	roleType, roleTypeErr := dbHelper.GetRoleByUserID(userID)
	if roleTypeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting role")
		return
	}

	fmt.Println(roleType)
	fmt.Println(userID)
	accessToken, accessErr := middleware.GenerateAccessToken(userID, roleType)
	refreshToken, refreshErr := middleware.GenerateRefreshToken(userID)
	if accessErr != nil || refreshErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate jwt token")
		return
	}

	_, sessErr := dbHelper.CreateSession(userID, refreshToken)
	if sessErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":       "user logged in successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	if EncodeErr != nil {
		return
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request")
		return
	}

	userID := middleware.UserContext(r)
	if err := dbHelper.Logout(userID, body.RefreshToken); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "logout failed")
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successfully",
	})
	if EncodeErr != nil {
		return
	}
}
