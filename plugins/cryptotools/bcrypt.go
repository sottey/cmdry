// Package cryptotools provides local password-hash utilities.
package cryptotools

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash creates a bcrypt hash using the requested cost.
func Hash(value string, cost int) (string, error) {
	if value == "" {
		return "", fmt.Errorf("text to hash is required")
	}
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(value), cost)
	if err != nil {
		return "", fmt.Errorf("create bcrypt hash: %w", err)
	}
	return string(hash), nil
}

// Verify reports whether value matches an existing bcrypt hash.
func Verify(value, hash string) (bool, error) {
	if value == "" || hash == "" {
		return false, fmt.Errorf("text and bcrypt hash are required")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	if err == nil {
		return true, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, fmt.Errorf("read bcrypt hash: %w", err)
}
