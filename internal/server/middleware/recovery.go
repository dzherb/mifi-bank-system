package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dzherb/mifi-bank-system/internal/server/responses"
)

type ctxKey string

const panicKey ctxKey = "panic_value"

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// It works but doesn't look nice...
				// We need to pass the panic value to the logging middleware
				ctx := context.WithValue(r.Context(), panicKey, rec)
				*r = *r.WithContext(ctx)

				w.WriteHeader(http.StatusInternalServerError)
				responses.WriteError(
					w,
					fmt.Errorf("internal server error"),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
