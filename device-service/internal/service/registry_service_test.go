package service

import (
	"context"
	"testing"

	"github.com/industrial-iot/device-service/internal/model"
)

func TestRegisterDevice_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     model.CreateDeviceRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "missing serial_number",
			req: model.CreateDeviceRequest{
				Name:      "No SN",
				FactoryID: "HN",
			},
			wantErr: true,
			errMsg:  "serial_number is required",
		},
		{
			name: "empty serial_number",
			req: model.CreateDeviceRequest{
				SerialNumber: "   ",
				Name:         "Empty SN",
				FactoryID:    "HN",
			},
			wantErr: true,
			errMsg:  "serial_number is required",
		},
		{
			name: "missing name",
			req: model.CreateDeviceRequest{
				SerialNumber: "SN-003",
				FactoryID:    "HN",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing factory_id",
			req: model.CreateDeviceRequest{
				SerialNumber: "SN-004",
				Name:         "No Factory",
			},
			wantErr: true,
			errMsg:  "factory_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &RegistryService{repo: nil, logger: nil}
			_, err := svc.RegisterDevice(context.Background(), &tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("error = %q, want %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestUpdateDeviceStatus_Validation(t *testing.T) {
	svc := &RegistryService{repo: nil, logger: nil}

	tests := []struct {
		status  model.DeviceStatus
		wantErr bool
	}{
		{model.DeviceStatus("invalid"), true},
		{model.DeviceStatus("unknown"), true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			err := svc.UpdateDeviceStatus(context.Background(), "fake-id", tt.status)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for status '%s' but got nil", tt.status)
				}
				if err.Error()[:8] != "invalid " {
					t.Errorf("error should start with 'invalid ', got: %s", err.Error())
				}
			}
		})
	}
}
