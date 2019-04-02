package util

import (
	"fmt"
	"math/rand"
	"time"
)

var chars = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("test123")
}

// GetRandString returnes a random stirng
// of the defined length.
func GetRandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
