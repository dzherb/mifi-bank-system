package server

import (
	"github.com/dzherb/mifi-bank-system/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(r *mux.Router) {
	r.Handle("/health", handlers.HealthHandler).Methods(http.MethodGet)
}
