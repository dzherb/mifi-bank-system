package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dzherb/mifi-bank-system/internal/pkg/responses"
	"github.com/dzherb/mifi-bank-system/internal/services/auth"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var authService = auth.NewService()

var RegisterHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		req := registerRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			responses.WriteError(w, err)

			return
		}

		resp, err := authService.Register(req.Email, req.Username, req.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			responses.WriteError(w, err)

			return
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responses.WriteError(w, err)
		}
	})

var LoginHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		req := loginRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			responses.WriteError(w, err)

			return
		}

		resp, err := authService.Login(req.Username, req.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			responses.WriteError(w, err)

			return
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responses.WriteError(w, err)
		}
	},
)
