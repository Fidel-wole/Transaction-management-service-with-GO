package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateAccountNumber generates a random 10-digit account number
func GenerateAccountNumber() string {
	// Create a new random generator with the current time as the seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a 10-digit number (random number between 1000000000 and 9999999999)
	accountNumber := r.Int63n(9000000000) + 1000000000

	// Convert the number to a string and return
	return fmt.Sprintf("%010d", accountNumber)
}
