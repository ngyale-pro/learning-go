package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// If the format specifier includes a %w verb with an error operand
		// the returned error will implement an Unwrap method returning the operand
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), err
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
