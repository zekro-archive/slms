package main

import (
	"crypto/sha256"
	"fmt"
)

func GetSHA256Hash(data string) string {
	return fmt.Sprintf("%x", sha256.New().Sum([]byte(data)))
}
