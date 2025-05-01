package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

var HealthHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			WriteErrorResponse(w, err)
		}
	},
)
