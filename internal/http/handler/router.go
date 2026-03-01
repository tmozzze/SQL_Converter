package handler

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterRoutes - register routes for http
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/upload", h.Upload)
	// swagger docs http://localhost:8080/swagger/index.html.
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
}
