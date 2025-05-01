package server

import (
	"net/http"

	"github.com/dzherb/mifi-bank-system/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.Handle("/health", handlers.HealthHandler).Methods(http.MethodGet)
}
