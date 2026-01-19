package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"smarthome/models"
)

// DeviceService sends device lifecycle events to the device microservice.
type DeviceService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewDeviceService creates a new device service client.
func NewDeviceService(baseURL string) *DeviceService {
	return &DeviceService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type deviceCreateRequest struct {
	DeviceID  string `json:"deviceId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Location  string `json:"location"`
	Status    string `json:"status"`
	Unit      string `json:"unit,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// RegisterDevice registers a device in the device service.
func (s *DeviceService) RegisterDevice(sensor models.Sensor) error {
	payload := deviceCreateRequest{
		DeviceID: fmt.Sprintf("%d", sensor.ID),
		Name:     sensor.Name,
		Type:     string(sensor.Type),
		Location: sensor.Location,
		Status:   sensor.Status,
		Unit:     sensor.Unit,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal device payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.BaseURL+"/devices", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build register request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("register device request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("register device failed: status %d", resp.StatusCode)
	}

	return nil
}
