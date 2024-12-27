package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	ua := NewApplication()

	r.Get("/users", HandleFindAll(ua))
	r.Get("/users/{id}", HandleFindById(ua))
	r.Post("/users", HandleInsert(ua))
	r.Put("/users/{id}", HandleUpdate(ua))
	r.Delete("/users/{id}", HandleDelete(ua))

	return r
}
