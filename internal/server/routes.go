package server

import (
	"net/http"

	handlers2 "github.com/dzherb/mifi-bank-system/internal/server/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.Handle("/health", handlers2.HealthHandler).Methods(http.MethodGet)
	r.Handle("/auth/register", handlers2.RegisterHandler).
		Methods(http.MethodPost)
	r.Handle("/auth/login", handlers2.LoginHandler).Methods(http.MethodPost)
	r.Handle("/users/me", handlers2.CurrentUserHandler).Methods(http.MethodGet)
}
