package middleware

import (
	"RMS/utils"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
	"time"
)

type ContextKey string

const (
	userContext ContextKey = "userKey"
)

type Claims struct {
	UserID string   `json:"user_id"`
	Role   []string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func GenerateAccessToken(userID string, roleType []string) (string, error) {
	accessClaims := &Claims{
		UserID: userID,
		Role:   roleType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessJWT.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
func GenerateRefreshToken(userID string) (string, error) {
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshJWT.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return refreshToken, nil

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ResponseError(w, http.StatusBadRequest, "invalid token")
			return
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		accessClaims := &Claims{}
		token, err := jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		// access token is valid
		if err == nil && token.Valid {
			ctx := context.WithValue(r.Context(), userContext, accessClaims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else { // access token invalid
			utils.ResponseError(w, http.StatusUnauthorized, "token expired")
			return
		}

	})
}

func UserContext(r *http.Request) string {
	if user, ok := r.Context().Value(userContext).(string); ok && user != "" {
		return user
	}
	return ""

}
