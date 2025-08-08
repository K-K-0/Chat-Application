package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Decode the base64-encoded JWT secret
	base64Secret := "/vBlkT0tyW8dYZ0lhXttcH5UY52Ayw9WSg+kFqNwftGIRWy5VeahsPlfhhmujjKbufdNUNJUISUM3Tqy3HR7FA=="
	jwtSecret, err := base64.StdEncoding.DecodeString(base64Secret)
	if err != nil {
		panic("Failed to decode JWT secret: " + err.Error())
	}

	claims := jwt.MapClaims{
		"sub":   "A@gmail.com",
		"role":  "authenticated",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"email": "A@gmail.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		panic(err)
	}

	fmt.Println("Generated JWT token:\n", signedToken)
}
