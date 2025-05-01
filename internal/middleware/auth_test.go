package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/middleware"
	"github.com/dzherb/mifi-bank-system/internal/security"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	userID := r.Context().Value(middleware.UserIDKey).(int)
	w.Write([]byte(strconv.Itoa(userID))) //nolint:errcheck
})

func initSecurity() {
	security.Init(&config.Config{
		SecretKey:      "secret",
		AccessTokenTTL: time.Hour,
	})
}

func TestAuthMiddlewareSuccess(t *testing.T) {
	initSecurity()

	userID := 25

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	token, err := security.IssueAccessToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	middleware.AuthRequired(handler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("got status %v, want %v", status, http.StatusOK)
	}

	if rr.Body.String() != strconv.Itoa(userID) {
		t.Errorf("got body %v, want %v", rr.Body.String(), strconv.Itoa(userID))
	}
}

func TestAuthMiddlewareFail(t *testing.T) {
	initSecurity()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer fake_token")

	rr := httptest.NewRecorder()

	middleware.AuthRequired(handler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("got status %v, want %v", status, http.StatusOK)
	}
}
