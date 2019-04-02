package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// CreateHash creates a hash from a string
// using the blowfish algorithm.
func CreateHash(s string, rounds int) (string, error) {
	bHash, err := bcrypt.GenerateFromPassword([]byte(s), rounds)
	return string(bHash), err
}

// CheckHash returns if the passed, bcrypt
// generated hash matches the passed string.
func CheckHash(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}
