package health

import (
	"net/http"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// HealthResponse represents the JSON response structure
type HealthResponse struct {
	Status string `json:"status"`
}

// Handler Handles health endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Info("Returning value for health endpoint")

	response := HealthResponse{
		Status: "ok",
	}

	// Convert the response struct to JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}