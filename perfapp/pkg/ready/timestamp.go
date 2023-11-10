package ready

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"

	"ocp.performance.io/perfapp/internal/perf"
	log "github.com/sirupsen/logrus"
)

// Tables Euler workload required tables
var Tables = map[string]string{"ts": "CREATE TABLE IF NOT EXISTS ts (date TIMESTAMP)"}

// TimestampResponse represents the JSON response structure for timestamp requests
type TimestampResponse struct {
	Status   string    `json:"status"`
	Query    string `json:"query"`
	Duration int64     `json:"duration_ns"`
}

// Handler Handle timestamp requests
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Info("Inserting timestamp record in table")
	now := time.Now()
	insert := fmt.Sprintf("INSERT INTO ts VALUES ('%s')", now.Format(time.RFC3339))
	if err := perf.QueryDB(insert); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, err.Error())
	} else {
		response := TimestampResponse{
			Status:   "ok",
			Query: insert,
			Duration: time.Since(now).Nanoseconds(),
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

		log.Printf("Timestamp inserted in %v ns", time.Since(now).Nanoseconds())
		perf.HTTPRequestDuration.Observe(time.Since(now).Seconds())
	}
}
