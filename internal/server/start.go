package server

import (
	"context"
	"net/http"
	"time"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/middleware"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const writeTimeout = 15 * time.Second
const readTimeout = 15 * time.Second

var srv *http.Server

func Start(cfg *config.Config) error {
	r := mux.NewRouter()
	RegisterRoutes(r)

	r.Use(
		middleware.LoggingMiddleware,
		middleware.JSONMiddleware,
		middleware.RecoveryMiddleware,
		middleware.RateLimiter(),
	)

	srv = &http.Server{
		Handler:      r,
		Addr:         cfg.ServerHost + ":" + cfg.ServerPort,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	log.Infof("starting the server on %s:%s", cfg.ServerHost, cfg.ServerPort)

	return srv.ListenAndServe()
}

func Shutdown(ctx context.Context) error {
	return srv.Shutdown(ctx)
}
