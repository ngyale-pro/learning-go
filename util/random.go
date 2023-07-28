package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and mex
func RandomInteger(min, max int64) int64 {
	return min + rand.Int63n(int64(max-min+1))
}

/* // RandomString generates a random string of length n
func RandomString(n int) string {
	var generated_string = ""
	for i := 0; i < n; i++ {
		generated_string += string(alphabet[rand.Intn(26)])
	}
	return generated_string
} */

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	var alphabet_length = len(alphabet)

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(alphabet_length)])
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInteger(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "YEN"}
	return currencies[rand.Intn(len(currencies))]
}
