package stream

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	hlsDir string
}

func (h *Handler) RegisterRoutes(r chi.Router) {

}
