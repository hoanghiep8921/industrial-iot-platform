package model

import "time"

// DeviceShadow stores the desired and reported state of a device
type DeviceShadow struct {
	DeviceID      string         `json:"device_id"`
	DesiredState  map[string]any `json:"desired_state"`
	ReportedState map[string]any `json:"reported_state"`
	Delta         map[string]any `json:"delta"`
	Version       int64          `json:"version"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// ShadowUpdate is the payload for updating desired or reported state
type ShadowUpdate struct {
	State map[string]any `json:"state"`
}

// DeltaResult is computed when desired and reported diverge
type DeltaResult struct {
	HasDelta bool
	Delta    map[string]any
}
