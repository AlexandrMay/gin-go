package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})
	tokenString, err := token.SignedString([]byte("supersecret"))
	if err != nil {
		return "", fmt.Errorf("Could not generate token: %w", err)
	}
	return tokenString, nil
}

func VerifyToken(token string) (int64, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("supersecret"), nil
	})
	if err != nil {
		return 0, fmt.Errorf("Could not parse token: %w", err)
	}

	isValid := parsed.Valid
	if !isValid {
		return 0, fmt.Errorf("Token is not valid")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("Could not parse claims")
	}

	userId := int64(claims["userId"].(float64))

	return userId, nil
}
