package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func NewBcryptHash(text string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

func CompareBcryptHash(hash string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false
	}

	return true
}
