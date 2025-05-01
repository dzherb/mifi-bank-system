package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
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
