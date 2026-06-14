package model

import (
	"time"
)

// DeviceStatus represents the current state of a device
type DeviceStatus string

const (
	StatusOnline      DeviceStatus = "online"
	StatusOffline     DeviceStatus = "offline"
	StatusMaintenance DeviceStatus = "maintenance"
	StatusError       DeviceStatus = "error"
)

// ValidDeviceStatuses is the list of allowed device status values
var ValidDeviceStatuses = []DeviceStatus{
	StatusOnline,
	StatusOffline,
	StatusMaintenance,
	StatusError,
}

// IsValid checks if the status value is allowed
func (s DeviceStatus) IsValid() bool {
	for _, v := range ValidDeviceStatuses {
		if s == v {
			return true
		}
	}
	return false
}

// Device represents a registered IoT device
type Device struct {
	ID              string         `json:"id"`
	SerialNumber    string         `json:"serial_number"`
	Name            string         `json:"name"`
	Model           string         `json:"model,omitempty"`
	Vendor          string         `json:"vendor,omitempty"`
	FactoryID       string         `json:"factory_id"`
	Area            string         `json:"area,omitempty"`
	Line            string         `json:"line,omitempty"`
	Protocol        string         `json:"protocol"`
	Capabilities    []string       `json:"capabilities"`
	Status          DeviceStatus   `json:"status"`
	FirmwareVersion string         `json:"firmware_version,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
	LastSeenAt      *time.Time     `json:"last_seen_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// CreateDeviceRequest is the payload for registering a new device
type CreateDeviceRequest struct {
	SerialNumber    string         `json:"serial_number"`              // required
	Name            string         `json:"name"`                       // required
	Model           string         `json:"model,omitempty"`
	Vendor          string         `json:"vendor,omitempty"`
	FactoryID       string         `json:"factory_id"`                 // required
	Area            string         `json:"area,omitempty"`
	Line            string         `json:"line,omitempty"`
	Protocol        string         `json:"protocol,omitempty"`
	Capabilities    []string       `json:"capabilities,omitempty"`
	FirmwareVersion string         `json:"firmware_version,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

// UpdateDeviceRequest is the payload for updating a device
type UpdateDeviceRequest struct {
	Name            *string         `json:"name,omitempty"`
	Model           *string         `json:"model,omitempty"`
	Vendor          *string         `json:"vendor,omitempty"`
	FactoryID       *string         `json:"factory_id,omitempty"`
	Area            *string         `json:"area,omitempty"`
	Line            *string         `json:"line,omitempty"`
	Protocol        *string         `json:"protocol,omitempty"`
	Capabilities    *[]string       `json:"capabilities,omitempty"`
	FirmwareVersion *string         `json:"firmware_version,omitempty"`
	Metadata        *map[string]any `json:"metadata,omitempty"`
}

// UpdateStatusRequest is the payload for changing device status
type UpdateStatusRequest struct {
	Status DeviceStatus `json:"status"`
}

// DeviceFilter holds query parameters for listing devices
type DeviceFilter struct {
	FactoryID string
	Area      string
	Status    DeviceStatus
	Protocol  string
	Search    string // search in name or serial_number
	Page      int
	Limit     int
}

// DeviceStats holds device statistics
type DeviceStats struct {
	Total     int                    `json:"total"`
	ByFactory map[string]int         `json:"by_factory"`
	ByStatus  map[DeviceStatus]int   `json:"by_status"`
	ByModel   map[string]int         `json:"by_model"`
}
