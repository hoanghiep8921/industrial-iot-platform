package model

import (
	"testing"
	"time"
)

func TestDeviceShadow(t *testing.T) {
	now := time.Now()
	shadow := DeviceShadow{
		DeviceID:      "dev-001",
		DesiredState:  map[string]any{"temp": 200},
		ReportedState: map[string]any{"temp": 195},
		Delta:         map[string]any{"temp": 200},
		Version:       1,
		UpdatedAt:     now,
	}

	if shadow.DeviceID != "dev-001" {
		t.Errorf("DeviceID = %s, want dev-001", shadow.DeviceID)
	}
	if shadow.Version != 1 {
		t.Errorf("Version = %d, want 1", shadow.Version)
	}
	if len(shadow.Delta) != 1 {
		t.Errorf("Delta length = %d, want 1", len(shadow.Delta))
	}
}

func TestShadowUpdate(t *testing.T) {
	update := ShadowUpdate{
		State: map[string]any{"temp": 200, "speed": 1500},
	}

	if len(update.State) != 2 {
		t.Errorf("State length = %d, want 2", len(update.State))
	}
}

func TestDeltaResult(t *testing.T) {
	dr := DeltaResult{
		HasDelta: true,
		Delta:    map[string]any{"temp": 200},
	}

	if !dr.HasDelta {
		t.Error("HasDelta should be true")
	}
	if len(dr.Delta) != 1 {
		t.Errorf("Delta length = %d, want 1", len(dr.Delta))
	}
}
