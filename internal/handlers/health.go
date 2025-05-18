package handlers

import (
	"encoding/json"
	"net/http"

	http2 "github.com/dzherb/mifi-bank-system/internal/pkg/responses"
)

type HealthResponse struct {
	Status string `json:"status"`
}

var HealthHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http2.WriteError(w, err)
		}
	},
)
