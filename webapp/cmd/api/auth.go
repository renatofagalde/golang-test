package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwtTokenExpiry = time.Minute * 15
const refreshTokenExpiry = time.Hour * 24

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func (app *application) validateToken(w http.ResponseWriter, request *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	authHeader := request.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}
}
