package middleware

import (
	"context"
	"fmt"
	"github.com/dzherb/mifi-bank-system/internal/handlers"
	"github.com/dzherb/mifi-bank-system/internal/security"
	"net/http"
	"strings"
)

type userIDCtxKey string

const UserIDKey userIDCtxKey = "userID"
const TokenPrefix = "Bearer "

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, TokenPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			handlers.WriteErrorResponse(w, fmt.Errorf("token must be provided"))
			return
		}

		token := strings.TrimPrefix(authHeader, TokenPrefix)
		userID, err := security.ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			handlers.WriteErrorResponse(w, err)
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userID))
		next.ServeHTTP(w, r)
	})
}
