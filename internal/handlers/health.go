package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HealthHandler struct{}

type HealthResponse struct {
	Status string `json:"status"`
}

func (hh HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	if err != nil {
		log.Error(err)
	}
}
