package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = os.Getenv("JWT_SECRET")

func GenerateToken(UserId uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = UserId
	claims["exp"] = time.Now().Add(time.Hour * 24)
	return token.SignedString([]byte(jwtKey))
}

func ValidateToken(tokenStr string) (uint, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("error parsing user ID from claims")
		}
		return uint(userIDFloat), nil
	}

	return 0, fmt.Errorf("invalid token")
}

func ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(jwtKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if idClaim, ok := claims["id"].(float64); ok {
			return uint(idClaim), nil
		}

		return 0, errors.New("token doesn't contain valid user id")
	}

	return 0, errors.New("invalid token")
}
