package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	isPasswordCorrect := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if isPasswordCorrect != nil && isPasswordCorrect == bcrypt.ErrMismatchedHashAndPassword {
		return isPasswordCorrect
	}

	return nil
}
