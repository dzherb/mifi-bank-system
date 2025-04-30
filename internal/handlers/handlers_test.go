package handlers

import (
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	cfg := config.Load()

	_, err := storage.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	code := storage.WithTempDB(func() int {
		return storage.WithMigratedDB(m.Run)
	})
	os.Exit(code)
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	HealthHandler{}.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %v, got %v",
			http.StatusOK, status)
	}

	expected := `{"status":"ok"}`
	got := strings.TrimSpace(rr.Body.String())

	if got != expected {
		t.Errorf("expected body %v, got %v",
			expected, rr.Body.String())
	}
}
