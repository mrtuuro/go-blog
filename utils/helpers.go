package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password []byte) string {

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 8)
	if err != nil {
		log.Fatal(err)
	}

	return string(hashedPassword)
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("login"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
