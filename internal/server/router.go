package server

import (
	"narwhal/internal/upload"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(uploadManager *upload.Manager) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		upload.NewHandler(uploadManager).RegisterRoutes(r)
	})

	return r
}
