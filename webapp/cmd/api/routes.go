package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Post("/auth", app.authenticate)
	mux.Post("/refresh", app.refresh)

	mux.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
		var payload = struct {
			Message string `json:"message"`
		}{
			Message: "hello world",
		}
		_ = app.writeJSON(writer, http.StatusOK, payload)
	})

	mux.Route("/users", func(mux chi.Router) {
		mux.Get("/", app.allUsers)
		mux.Get("/{userID}", app.getUser)
		mux.Delete("/{userID}", app.getUser)
		mux.Put("/", app.insertUser)
		mux.Patch("/", app.updateUser)

	})

	return mux
}
