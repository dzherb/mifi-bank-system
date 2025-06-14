package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dzherb/mifi-bank-system/internal/security"
	"github.com/dzherb/mifi-bank-system/internal/server/responses"
)

type userIDCtxKey string

const UserIDKey userIDCtxKey = "userID"
const TokenPrefix = "Bearer "

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, TokenPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			responses.WriteError(w, fmt.Errorf("token must be provided"))

			return
		}

		token := strings.TrimPrefix(authHeader, TokenPrefix)

		userID, err := security.ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			responses.WriteError(w, err)
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userID))
		next.ServeHTTP(w, r)
	})
}
