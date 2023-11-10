package euler

import (
	"fmt"
	"math"
	"net/http"
	"time"
	"encoding/json"

	"ocp.performance.io/perfapp/internal/perf"
	log "github.com/sirupsen/logrus"
)

// Tables Euler workload required tables
var Tables = map[string]string{"euler": "CREATE TABLE IF NOT EXISTS euler (date TIMESTAMP, elapsed FLOAT(24))"}

// Defines response type
type Response struct {
	Status  string  `json:"status"`
	Query	string `json:"query"`
	Message string  `json:"message,omitempty"`
	Duration float64 `json:"duration,omitempty"`
}

// Handler Handle requests to compute euler number aproximation
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Computing euler approximation")
	now := time.Now()
	calcEuler()
	duration := time.Since(now).Seconds()
	insert := fmt.Sprintf("INSERT INTO euler VALUES ('%s', '%f')", now.Format(time.RFC3339), duration)

	response := Response{} // Initialize a Response struct

	if err := perf.QueryDB(insert); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		response.Status = "error"
		response.Query = insert
		response.Message = err.Error()
	} else {
		w.WriteHeader(http.StatusOK)
		response.Status = "success"
		response.Query = insert
		response.Duration = duration
	}

	// Convert the Response struct to JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	// Set the Content-Type header to indicate JSON response
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the HTTP response writer
	w.Write(jsonData)

	log.Printf("Euler approximation computed in %f seconds", duration)
	perf.HTTPRequestDuration.Observe(duration)
}

func calcEuler() {
	var n float64
	var x float64
	for math.E > x {
		x = math.Pow((1 + 1/n), n)
		n++
	}
}
