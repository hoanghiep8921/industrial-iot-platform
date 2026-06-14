package model

import (
	"testing"
)

func TestDeviceStatus_IsValid(t *testing.T) {
	tests := []struct {
		status   DeviceStatus
		expected bool
	}{
		{StatusOnline, true},
		{StatusOffline, true},
		{StatusMaintenance, true},
		{StatusError, true},
		{DeviceStatus("unknown"), false},
		{DeviceStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.expected {
				t.Errorf("DeviceStatus.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidDeviceStatuses(t *testing.T) {
	if len(ValidDeviceStatuses) != 4 {
		t.Errorf("Expected 4 valid statuses, got %d", len(ValidDeviceStatuses))
	}

	expected := []DeviceStatus{StatusOnline, StatusOffline, StatusMaintenance, StatusError}
	for i, v := range expected {
		if ValidDeviceStatuses[i] != v {
			t.Errorf("ValidDeviceStatuses[%d] = %v, want %v", i, ValidDeviceStatuses[i], v)
		}
	}
}

func TestCreateDeviceRequest(t *testing.T) {
	req := CreateDeviceRequest{
		SerialNumber: "SN-TEST-001",
		Name:         "Test Device",
		FactoryID:    "HN",
		Protocol:     "mqtt",
		Capabilities: []string{"temperature", "speed"},
	}

	if req.SerialNumber != "SN-TEST-001" {
		t.Errorf("SerialNumber = %s, want SN-TEST-001", req.SerialNumber)
	}
	if req.Protocol != "mqtt" {
		t.Errorf("Protocol = %s, want mqtt", req.Protocol)
	}
	if len(req.Capabilities) != 2 {
		t.Errorf("Capabilities length = %d, want 2", len(req.Capabilities))
	}
}

func TestUpdateDeviceRequest(t *testing.T) {
	newName := "Updated Name"
	req := UpdateDeviceRequest{
		Name: &newName,
	}

	if req.Name == nil || *req.Name != "Updated Name" {
		t.Errorf("Name = %v, want 'Updated Name'", req.Name)
	}
}

func TestDeviceFilter(t *testing.T) {
	filter := DeviceFilter{
		FactoryID: "HN",
		Status:    StatusOnline,
		Page:      1,
		Limit:     10,
	}

	if filter.FactoryID != "HN" {
		t.Errorf("FactoryID = %s, want HN", filter.FactoryID)
	}
	if filter.Page != 1 {
		t.Errorf("Page = %d, want 1", filter.Page)
	}
}
