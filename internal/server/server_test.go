package server_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/server"
)

func TestServerStartAndShutdown(t *testing.T) {
	go func() {
		err := server.Start(&config.Config{
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

	// make sure the server has started
	time.Sleep(1 * time.Second)

	defer func(ctx context.Context) {
		err := server.Shutdown(ctx)
		if err != nil {
			t.Error(err)
		}
	}(t.Context())

	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
