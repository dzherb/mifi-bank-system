package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	res := ErrorResponse{Error: err.Error()}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.WithField("err", err).Error("Failed to write response")
	}
}
