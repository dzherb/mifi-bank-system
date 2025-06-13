package server

import (
	"net/http"

	"github.com/dzherb/mifi-bank-system/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.Handle("/health", handlers.HealthHandler).Methods(http.MethodGet)
	r.Handle("/auth/register", handlers.RegisterHandler).
		Methods(http.MethodPost)
	r.Handle("/auth/login", handlers.LoginHandler).Methods(http.MethodPost)
	r.Handle("/users/me", handlers.CurrentUserHandler).Methods(http.MethodGet)
}
