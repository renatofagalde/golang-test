package main

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Credenials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	var creds Credenials

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	user, err := app.DB.GetUserByEmail(creds.Username)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	tokenPairs, err := app.createTokenPair(user)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// ---------------------------------------------------
	// 🔥 SE HTTPS => usa cookie
	// 🔥 SE HTTP  => devolve refresh_token no JSON, igual antes
	// ---------------------------------------------------

	isHTTPS := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"

	if isHTTPS {
		// usar cookie apenas se HTTPS
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokenPairs.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   60 * 60 * 24 * 7, // 7 dias
		})

		// apenas o access token no body
		_ = app.writeJSON(w, http.StatusCreated, map[string]string{
			"access_token": tokenPairs.Token,
		})

		return
	}

	// ---------------------------------------------------
	// 🌐 HTTP → mantém comportamento antigo
	// ---------------------------------------------------
	_ = app.writeJSON(w, http.StatusCreated, tokenPairs)
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {

}

func (app *application) allUsers(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) insertUser(w http.ResponseWriter, r *http.Request) {

}
