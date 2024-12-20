package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash, %w", err)
	}

	return string(bytes), nil
}

func CheckPwdHash(pwd, hashedPwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	if err != nil {
		return false, fmt.Errorf("wrong password: %w", err)
	}

	return true, nil
}
