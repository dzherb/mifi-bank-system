package server

import (
	"context"
	"errors"
	"github.com/dzherb/mifi-bank-system/internal/config"
	"net/http"
	"testing"
)

func TestServerStartAndShutdown(t *testing.T) {
	go func() {
		err := Start(&config.Config{
			ServerHost: "localhost",
			ServerPort: "8080",
		})

		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			t.Error(err)
		}
	}()

	defer func(ctx context.Context) {
		err := Shutdown(ctx)
		if err != nil {
			t.Error(err)
		}
	}(context.Background())

	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
