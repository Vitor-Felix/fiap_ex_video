package application

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "fiap-x-super-secret"

func GenerateJWT(userID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
