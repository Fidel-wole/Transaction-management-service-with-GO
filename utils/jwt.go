package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

// GenerateToken generates a JWT token with email and user ID claims.
func GenerateToken(email string, userId int64) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Set expiration time
	}
	fmt.Printf("Claims: %v\n", claims)
	// Create the token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return the token
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) (int64, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return 0, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, errors.New("invalid claims")
    }

    // Debugging: print the claims to ensure they are correct
    fmt.Printf("Claims: %v\n", claims)

    // Extract the userId from the claims and convert it to int64
    userIdFloat, ok := claims["userId"].(float64)
    if !ok {
        return 0, errors.New("userId is not a valid number")
    }
    userId := int64(userIdFloat)

    // Debugging: print the userId to check its value
    fmt.Printf("Extracted userId: %d\n", userId)

    return userId, nil
}
