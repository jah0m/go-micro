package main

import (
	"authentication/data"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateAccessToken creates a new JWT token for the user
func (app *Config) generateAccessToken(user *data.User) (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Sign the token with the provided key
	tokenString, err := token.SignedString([]byte(app.ACCESS_KEY))
	if err != nil {
		return "", errors.New("failed to generate JWT")
	}

	return tokenString, nil
}

// generateRefreshToken creates a new JWT token for the user
func (app *Config) generateRefreshToken(user *data.User) (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	// Sign the token with the provided key
	tokenString, err := token.SignedString([]byte(app.REFRESH_KEY))
	if err != nil {
		return "", errors.New("failed to generate JWT")
	}

	return tokenString, nil
}
