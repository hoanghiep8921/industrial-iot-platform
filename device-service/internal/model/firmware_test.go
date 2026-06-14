package model

import (
	"testing"
)

func TestUpdateStatus_IsValid(t *testing.T) {
	tests := []struct {
		status   UpdateStatus
		expected bool
	}{
		{UpdatePending, true},
		{UpdateDownloading, true},
		{UpdateInstalling, true},
		{UpdateSuccess, true},
		{UpdateFailed, true},
		{UpdateRolledBack, true},
		{UpdateStatus("unknown"), false},
		{UpdateStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.expected {
				t.Errorf("UpdateStatus.IsValid(%s) = %v, want %v", tt.status, got, tt.expected)
			}
		})
	}
}

func TestUpdateStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		status   UpdateStatus
		expected bool
	}{
		{UpdateSuccess, true},
		{UpdateFailed, true},
		{UpdateRolledBack, true},
		{UpdatePending, false},
		{UpdateDownloading, false},
		{UpdateInstalling, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := tt.status.IsTerminal()
			if got != tt.expected {
				t.Errorf("UpdateStatus.IsTerminal(%s) = %v, want %v", tt.status, got, tt.expected)
			}
		})
	}
}

func TestFirmwareStatusConstants(t *testing.T) {
	if FirmwareDraft != "draft" {
		t.Errorf("FirmwareDraft = %s, want draft", FirmwareDraft)
	}
	if FirmwareReleased != "released" {
		t.Errorf("FirmwareReleased = %s, want released", FirmwareReleased)
	}
	if FirmwareDeprecated != "deprecated" {
		t.Errorf("FirmwareDeprecated = %s, want deprecated", FirmwareDeprecated)
	}
}

func TestFirmware_Helpers(t *testing.T) {
	fw := Firmware{}

	// Nil values return empty/zero
	if fw.FileNameStr() != "" {
		t.Errorf("FileNameStr() = %s, want empty", fw.FileNameStr())
	}
	if fw.FileSizeVal() != 0 {
		t.Errorf("FileSizeVal() = %d, want 0", fw.FileSizeVal())
	}
	if fw.ChecksumStr() != "" {
		t.Errorf("ChecksumStr() = %s, want empty", fw.ChecksumStr())
	}
	if fw.MinioPathStr() != "" {
		t.Errorf("MinioPathStr() = %s, want empty", fw.MinioPathStr())
	}

	// Non-nil values
	name := "test.bin"
	checksum := "abc123"
	path := "/path/to/file"
	var size int64 = 1024

	fw.FileName = &name
	fw.ChecksumSHA256 = &checksum
	fw.MinioPath = &path
	fw.FileSize = &size

	if fw.FileNameStr() != "test.bin" {
		t.Errorf("FileNameStr() = %s, want test.bin", fw.FileNameStr())
	}
	if fw.FileSizeVal() != 1024 {
		t.Errorf("FileSizeVal() = %d, want 1024", fw.FileSizeVal())
	}
	if fw.ChecksumStr() != "abc123" {
		t.Errorf("ChecksumStr() = %s, want abc123", fw.ChecksumStr())
	}
	if fw.MinioPathStr() != "/path/to/file" {
		t.Errorf("MinioPathStr() = %s, want /path/to/file", fw.MinioPathStr())
	}
}

func TestFirmwareUpdate_Helpers(t *testing.T) {
	u := FirmwareUpdate{}

	if u.CampaignNameStr() != "" {
		t.Errorf("CampaignNameStr() = %s, want empty", u.CampaignNameStr())
	}
	if u.FromVersionStr() != "" {
		t.Errorf("FromVersionStr() = %s, want empty", u.FromVersionStr())
	}
	if u.ErrorMsg() != "" {
		t.Errorf("ErrorMsg() = %s, want empty", u.ErrorMsg())
	}

	campaign := "test-campaign"
	from := "v1.0.0"
	errMsg := "connection lost"

	u.CampaignName = &campaign
	u.FromVersion = &from
	u.ErrorMessage = &errMsg

	if u.CampaignNameStr() != "test-campaign" {
		t.Errorf("CampaignNameStr() = %s, want test-campaign", u.CampaignNameStr())
	}
	if u.FromVersionStr() != "v1.0.0" {
		t.Errorf("FromVersionStr() = %s, want v1.0.0", u.FromVersionStr())
	}
	if u.ErrorMsg() != "connection lost" {
		t.Errorf("ErrorMsg() = %s, want connection lost", u.ErrorMsg())
	}
}

func TestCreateFirmwareRequest(t *testing.T) {
	req := CreateFirmwareRequest{
		Version:     "v1.0.0",
		Model:       "CNC-L500",
		Description: "Test firmware",
	}

	if req.Version != "v1.0.0" {
		t.Errorf("Version = %s, want v1.0.0", req.Version)
	}
}

func TestCreateUpdateRequest(t *testing.T) {
	req := CreateUpdateRequest{
		FirmwareID: "fw-123",
		DeviceIDs:  []string{"dev-1", "dev-2", "dev-3", "dev-4", "dev-5"},
		Priority:   1,
		MaxRetries: 3,
	}

	if req.FirmwareID != "fw-123" {
		t.Errorf("FirmwareID = %s, want fw-123", req.FirmwareID)
	}
	if len(req.DeviceIDs) != 5 {
		t.Errorf("DeviceIDs length = %d, want 5", len(req.DeviceIDs))
	}
}
