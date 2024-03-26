package controller

import (
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

func generateResponseError(status int, message string) Response {
	return Response{
		Status:  status,
		Message: message,
	}
}

func generateResponseSuccess(status int, message string, data any) Response {
	return Response{
		Status:  status,
		Data:    data,
		Message: message,
	}
}

func createToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token kedaluwarsa dalam 24 jam
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
