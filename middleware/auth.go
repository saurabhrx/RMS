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
	roleContext ContextKey = "roleKey"
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
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
			ctx = context.WithValue(ctx, roleContext, accessClaims.Role)
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

func RoleContext(r *http.Request) []string {
	if roles, ok := r.Context().Value(roleContext).([]string); ok {
		return roles
	}
	return []string{}
}

func AuthRole(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{})
	for _, role := range allowedRoles {
		roleSet[strings.ToLower(role)] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles := RoleContext(r)
			for _, role := range userRoles {
				if _, ok := roleSet[strings.ToLower(role)]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.ResponseError(w, http.StatusForbidden, "unauthorized role")
		})
	}
}
