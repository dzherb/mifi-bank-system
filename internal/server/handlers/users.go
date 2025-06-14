package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/server/middleware"
	"github.com/dzherb/mifi-bank-system/internal/server/responses"
)

var CurrentUserHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)

		user, err := repo.NewUserRepository().Get(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responses.WriteError(w, fmt.Errorf("failed to get user"))
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responses.WriteError(w, err)
		}
	})
