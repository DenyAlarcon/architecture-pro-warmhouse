package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type temperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/temperature", handleTemperatureQuery)
	http.HandleFunc("/temperature/", handleTemperatureByID)

	port := getEnv("PORT", "8081")
	addr := ":" + port
	log.Printf("temperature-api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleTemperatureQuery(w http.ResponseWriter, r *http.Request) {
	location := strings.TrimSpace(r.URL.Query().Get("location"))
	sensorID := strings.TrimSpace(r.URL.Query().Get("sensorId"))
	if sensorID == "" {
		sensorID = strings.TrimSpace(r.URL.Query().Get("sensorID"))
	}

	location, sensorID = normalizeLocationAndSensorID(location, sensorID)
	writeTemperatureResponse(w, location, sensorID)
}

func handleTemperatureByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/temperature/")
	sensorID := strings.TrimSpace(path)
	location := strings.TrimSpace(r.URL.Query().Get("location"))

	location, sensorID = normalizeLocationAndSensorID(location, sensorID)
	writeTemperatureResponse(w, location, sensorID)
}

func normalizeLocationAndSensorID(location, sensorID string) (string, string) {
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	return location, sensorID
}

func writeTemperatureResponse(w http.ResponseWriter, location, sensorID string) {
	value := 18.0 + rand.Float64()*8.0
	value = math.Round(value*10) / 10

	resp := temperatureResponse{
		Value:       value,
		Unit:        "C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "ok",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Random temperature",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
