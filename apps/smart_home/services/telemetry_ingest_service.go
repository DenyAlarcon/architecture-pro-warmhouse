package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TelemetryIngestService sends telemetry records to the telemetry microservice.
type TelemetryIngestService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewTelemetryIngestService creates a new telemetry service client.
func NewTelemetryIngestService(baseURL string) *TelemetryIngestService {
	return &TelemetryIngestService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type telemetryRecord struct {
	DeviceID   string    `json:"deviceId"`
	Metric     string    `json:"metric"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit,omitempty"`
	RecordedAt time.Time `json:"recordedAt"`
}

// Record sends a telemetry record.
func (s *TelemetryIngestService) Record(deviceID, metric string, value float64, unit string, recordedAt time.Time) error {
	payload := telemetryRecord{
		DeviceID:   deviceID,
		Metric:     metric,
		Value:      value,
		Unit:       unit,
		RecordedAt: recordedAt,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal telemetry payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.BaseURL+"/telemetry", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build telemetry request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("telemetry request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telemetry request failed: status %d", resp.StatusCode)
	}

	return nil
}
