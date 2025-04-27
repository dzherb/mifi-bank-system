package server

import (
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/middleware"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Start(cfg *config.Config) error {
	r := mux.NewRouter()
	RegisterRoutes(r)

	r.Use(middleware.JSONMiddleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         cfg.ServerHost + ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("starting the server on %s:%s", cfg.ServerHost, cfg.ServerPort)
	return srv.ListenAndServe()
}
