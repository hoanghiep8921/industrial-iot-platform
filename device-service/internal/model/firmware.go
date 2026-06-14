package model

import "time"

// FirmwareStatus represents the lifecycle of a firmware
type FirmwareStatus string

const (
	FirmwareDraft      FirmwareStatus = "draft"
	FirmwareReleased   FirmwareStatus = "released"
	FirmwareDeprecated FirmwareStatus = "deprecated"
)

// Firmware represents a firmware binary registered in the system
type Firmware struct {
	ID             string         `json:"id"`
	Version        string         `json:"version"`
	Model          *string        `json:"model,omitempty"`
	Description    *string        `json:"description,omitempty"`
	FileName       *string        `json:"file_name,omitempty"`
	FileSize       *int64         `json:"file_size,omitempty"`
	ChecksumSHA256 *string        `json:"checksum_sha256,omitempty"`
	MinioPath      *string        `json:"-"`
	Status         FirmwareStatus `json:"status"`
	CreatedBy      *string        `json:"created_by,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	ReleasedAt     *time.Time     `json:"released_at,omitempty"`
}

// FileNameStr returns the file name or empty string
func (f *Firmware) FileNameStr() string {
	if f.FileName == nil { return "" }
	return *f.FileName
}

// FileSizeVal returns the file size or 0
func (f *Firmware) FileSizeVal() int64 {
	if f.FileSize == nil { return 0 }
	return *f.FileSize
}

// ChecksumStr returns the checksum or empty string
func (f *Firmware) ChecksumStr() string {
	if f.ChecksumSHA256 == nil { return "" }
	return *f.ChecksumSHA256
}

// MinioPathStr returns the minio path or empty string
func (f *Firmware) MinioPathStr() string {
	if f.MinioPath == nil { return "" }
	return *f.MinioPath
}

// CreateFirmwareRequest is the payload for registering firmware metadata
type CreateFirmwareRequest struct {
	Version     string `json:"version"`
	Model       string `json:"model,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedBy   string `json:"created_by,omitempty"`
}

// UpdateStatus represents the state of a firmware update on a device
type UpdateStatus string

const (
	UpdatePending     UpdateStatus = "pending"
	UpdateDownloading UpdateStatus = "downloading"
	UpdateInstalling  UpdateStatus = "installing"
	UpdateSuccess     UpdateStatus = "success"
	UpdateFailed      UpdateStatus = "failed"
	UpdateRolledBack  UpdateStatus = "rolled_back"
)

// ValidUpdateStatuses lists allowed transition states
var ValidUpdateStatuses = []UpdateStatus{
	UpdatePending,
	UpdateDownloading,
	UpdateInstalling,
	UpdateSuccess,
	UpdateFailed,
	UpdateRolledBack,
}

// IsValid checks if the update status value is allowed
func (s UpdateStatus) IsValid() bool {
	for _, v := range ValidUpdateStatuses {
		if s == v {
			return true
		}
	}
	return false
}

// IsTerminal returns true if the status is a final state
func (s UpdateStatus) IsTerminal() bool {
	return s == UpdateSuccess || s == UpdateFailed || s == UpdateRolledBack
}

// FirmwareUpdate represents an OTA update task for a device
type FirmwareUpdate struct {
	ID           string       `json:"id"`
	FirmwareID   string       `json:"firmware_id"`
	DeviceID     string       `json:"device_id"`
	CampaignName *string      `json:"campaign_name,omitempty"`
	Status       UpdateStatus `json:"status"`
	FromVersion  *string      `json:"from_version,omitempty"`
	ToVersion    string       `json:"to_version"`
	Priority     int          `json:"priority"`
	RetryCount   int          `json:"retry_count"`
	MaxRetries   int          `json:"max_retries"`
	StartedAt    *time.Time   `json:"started_at,omitempty"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty"`
	ErrorMessage *string      `json:"error_message,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (u *FirmwareUpdate) CampaignNameStr() string {
	if u.CampaignName == nil { return "" }
	return *u.CampaignName
}
func (u *FirmwareUpdate) FromVersionStr() string {
	if u.FromVersion == nil { return "" }
	return *u.FromVersion
}
func (u *FirmwareUpdate) ErrorMsg() string {
	if u.ErrorMessage == nil { return "" }
	return *u.ErrorMessage
}

// CreateUpdateRequest is the payload for creating a firmware update
type CreateUpdateRequest struct {
	FirmwareID   string   `json:"firmware_id"`
	DeviceID     string   `json:"device_id,omitempty"`   // single device
	DeviceIDs    []string `json:"device_ids,omitempty"`   // campaign (batch)
	CampaignName string   `json:"campaign_name,omitempty"`
	Priority     int      `json:"priority,omitempty"`
	MaxRetries   int      `json:"max_retries,omitempty"`
}

// UpdateStatusRequest is the payload for reporting update progress
type UpdateProgressRequest struct {
	Status       UpdateStatus `json:"status"`
	ErrorMessage string       `json:"error_message,omitempty"`
}
