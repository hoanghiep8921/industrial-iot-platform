package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/industrial-iot/device-service/internal/model"
)

func getFirmwareTestRepo(t *testing.T) *FirmwareRepo {
	t.Helper()
	repo := getTestRepo(t)
	return NewFirmwareRepo(repo)
}

func TestFirmwareRepo_CreateAndGetFirmware(t *testing.T) {
	r := getFirmwareTestRepo(t)
	ctx := context.Background()

	req := &model.CreateFirmwareRequest{
		Version:     "ITEST-T1781408501-v1.0.0",
		Model:       "CNC-L500",
		Description: "Integration test firmware",
		CreatedBy:   "test",
	}

	fw, err := r.CreateFirmware(ctx, req)
	if err != nil {
		t.Fatalf("CreateFirmware failed: %v", err)
	}
	if fw.ID == "" {
		t.Error("firmware ID should not be empty")
	}
	if fw.Status != model.FirmwareDraft {
		t.Errorf("status = %s, want draft", fw.Status)
	}

	// Get
	got, err := r.GetFirmwareByID(ctx, fw.ID)
	if err != nil {
		t.Fatalf("GetFirmwareByID failed: %v", err)
	}
	if got == nil {
		t.Fatal("firmware should exist")
	}
	if got.Version != "ITEST-T1781408501-v1.0.0" {
		t.Errorf("version = %s, want ITEST-T1781408501-v1.0.0", got.Version)
	}

	// Cleanup
	r.DeleteFirmware(ctx, fw.ID)
}

func TestFirmwareRepo_ListFirmwares(t *testing.T) {
	r := getFirmwareTestRepo(t)
	ctx := context.Background()

	ids := []string{}
	for i := 1; i <= 2; i++ {
		req := &model.CreateFirmwareRequest{
			Version: fmt.Sprintf("ITEST-T1781408501-list-v%d.0.0", i),
			Model:   "CNC-L500",
		}
		fw, err := r.CreateFirmware(ctx, req)
		if err != nil {
			t.Fatalf("CreateFirmware %d failed: %v", i, err)
		}
		ids = append(ids, fw.ID)
	}
	defer func() {
		for _, id := range ids {
			r.DeleteFirmware(ctx, id)
		}
	}()

	fws, err := r.ListFirmwares(ctx, "", "")
	if err != nil {
		t.Fatalf("ListFirmwares failed: %v", err)
	}
	if len(fws) < 2 {
		t.Errorf("expected at least 2 firmwares, got %d", len(fws))
	}

	// Filter by status
	drafts, err := r.ListFirmwares(ctx, "draft", "")
	if err != nil {
		t.Fatalf("ListFirmwares(draft) failed: %v", err)
	}
	for _, fw := range drafts {
		if fw.Status != model.FirmwareDraft {
			t.Errorf("expected draft status, got %s", fw.Status)
		}
	}
}

func TestFirmwareRepo_StatusLifecycle(t *testing.T) {
	r := getFirmwareTestRepo(t)
	ctx := context.Background()

	req := &model.CreateFirmwareRequest{
		Version: "ITEST-T1781408501-lifecycle-v1.0.0",
		Model:   "CNC-L500",
	}
	fw, err := r.CreateFirmware(ctx, req)
	if err != nil {
		t.Fatalf("CreateFirmware failed: %v", err)
	}
	defer r.DeleteFirmware(ctx, fw.ID)

	// Can't release without file
	err = r.UpdateFirmwareStatus(ctx, fw.ID, model.FirmwareReleased)
	// Actually this is allowed at repo level (service does the validation)

	// Update file
	err = r.UpdateFirmwareFile(ctx, fw.ID, "test.bin", 100, "abc123", "v1/test.bin")
	if err != nil {
		t.Fatalf("UpdateFirmwareFile failed: %v", err)
	}

	got, _ := r.GetFirmwareByID(ctx, fw.ID)
	if got.FileNameStr() != "test.bin" {
		t.Errorf("file_name = %s, want test.bin", got.FileNameStr())
	}
	if got.FileSizeVal() != 100 {
		t.Errorf("file_size = %d, want 100", got.FileSizeVal())
	}

	// Release
	err = r.UpdateFirmwareStatus(ctx, fw.ID, model.FirmwareReleased)
	if err != nil {
		t.Fatalf("UpdateFirmwareStatus(released) failed: %v", err)
	}

	got2, _ := r.GetFirmwareByID(ctx, fw.ID)
	if got2.Status != model.FirmwareReleased {
		t.Errorf("status = %s, want released", got2.Status)
	}
	if got2.ReleasedAt == nil {
		t.Error("released_at should be set")
	}

	// Deprecate
	err = r.UpdateFirmwareStatus(ctx, fw.ID, model.FirmwareDeprecated)
	if err != nil {
		t.Fatalf("UpdateFirmwareStatus(deprecated) failed: %v", err)
	}

	got3, _ := r.GetFirmwareByID(ctx, fw.ID)
	if got3.Status != model.FirmwareDeprecated {
		t.Errorf("status = %s, want deprecated", got3.Status)
	}
}

func TestFirmwareRepo_FirmwareUpdate(t *testing.T) {
	r := getFirmwareTestRepo(t)
	ctx := context.Background()
	repo := r.repo

	// Create a device
	d, err := repo.CreateDevice(ctx, &model.CreateDeviceRequest{
		SerialNumber: "ITEST-T1781408581-FWUPD",
		Name:         "FW Update Test",
		FactoryID:    "HN",
	})
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}
	defer repo.DeleteDevice(ctx, d.ID)

	// Create firmware
	fw, err := r.CreateFirmware(ctx, &model.CreateFirmwareRequest{
		Version: "ITEST-T1781408501-upd-v1.0.0",
	})
	if err != nil {
		t.Fatalf("CreateFirmware failed: %v", err)
	}
	defer r.DeleteFirmware(ctx, fw.ID)

	// Create update
	upd, err := r.CreateUpdate(ctx, fw.ID, d.ID, "test-campaign", "", fw.Version, 1, 2)
	if err != nil {
		t.Fatalf("CreateUpdate failed: %v", err)
	}
	if upd.Status != model.UpdatePending {
		t.Errorf("status = %s, want pending", upd.Status)
	}
	if upd.ToVersion != fw.Version {
		t.Errorf("to_version = %s, want %s", upd.ToVersion, fw.Version)
	}

	// Progress
	err = r.UpdateProgress(ctx, upd.ID, model.UpdateDownloading, "")
	if err != nil {
		t.Fatalf("UpdateProgress(downloading) failed: %v", err)
	}
	got, _ := r.GetUpdateByID(ctx, upd.ID)
	if got.Status != model.UpdateDownloading {
		t.Errorf("status = %s, want downloading", got.Status)
	}

	err = r.UpdateProgress(ctx, upd.ID, model.UpdateSuccess, "")
	if err != nil {
		t.Fatalf("UpdateProgress(success) failed: %v", err)
	}

	// Verify device firmware_version updated
	d2, _ := repo.GetDeviceByID(ctx, d.ID)
	if d2.FirmwareVersion != fw.Version {
		t.Errorf("device firmware_version = %s, want %s", d2.FirmwareVersion, fw.Version)
	}

	// List updates
	updates, err := r.ListUpdates(ctx, d.ID, "", "")
	if err != nil {
		t.Fatalf("ListUpdates failed: %v", err)
	}
	if len(updates) < 1 {
		t.Errorf("expected at least 1 update, got %d", len(updates))
	}

	// Rollback
	rollback, err := r.RollbackUpdate(ctx, upd.ID)
	if err != nil {
		t.Fatalf("RollbackUpdate failed: %v", err)
	}
	if rollback.ToVersion != "" {
		t.Errorf("rollback to_version = %s, want empty", rollback.ToVersion)
	}
}

func TestFirmwareRepo_RetryUpdate(t *testing.T) {
	r := getFirmwareTestRepo(t)
	ctx := context.Background()
	repo := r.repo

	d, err := repo.CreateDevice(ctx, &model.CreateDeviceRequest{
		SerialNumber: "ITEST-T1781408581-RETRY",
		Name:         "Retry Test",
		FactoryID:    "HN",
	})
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}
	defer repo.DeleteDevice(ctx, d.ID)

	fw, err := r.CreateFirmware(ctx, &model.CreateFirmwareRequest{Version: "ITEST-T1781408501-retry-v1.0.0"})
	if err != nil {
		t.Fatalf("CreateFirmware failed: %v", err)
	}
	defer r.DeleteFirmware(ctx, fw.ID)

	upd, _ := r.CreateUpdate(ctx, fw.ID, d.ID, "", "", fw.Version, 0, 2)

	// Mark as failed
	r.UpdateProgress(ctx, upd.ID, model.UpdateFailed, "error")

	// Retry
	err = r.RetryUpdate(ctx, upd.ID)
	if err != nil {
		t.Fatalf("RetryUpdate failed: %v", err)
	}

	got, _ := r.GetUpdateByID(ctx, upd.ID)
	if got.Status != model.UpdatePending {
		t.Errorf("status = %s, want pending after retry", got.Status)
	}
	if got.RetryCount != 1 {
		t.Errorf("retry_count = %d, want 1", got.RetryCount)
	}

	// Exhaust retries
	r.UpdateProgress(ctx, upd.ID, model.UpdateFailed, "error2")
	r.RetryUpdate(ctx, upd.ID)
	r.UpdateProgress(ctx, upd.ID, model.UpdateFailed, "error3")

	err = r.RetryUpdate(ctx, upd.ID)
	if err == nil {
		t.Error("expected error when max retries exceeded")
	}
}

func TestFirmwareRepo_BatchCreateUpdates(t *testing.T) {
	r := getFirmwareTestRepo(t)
	ctx := context.Background()
	repo := r.repo

	// Create devices
	deviceIDs := []string{}
	for i := 1; i <= 3; i++ {
		d, err := repo.CreateDevice(ctx, &model.CreateDeviceRequest{
			SerialNumber: fmt.Sprintf("ITEST-T1781408581-BATCH-%d", i),
			Name:         fmt.Sprintf("Batch Device %d", i),
			FactoryID:    "HN",
		})
		if err != nil {
			t.Fatalf("CreateDevice %d failed: %v", i, err)
		}
		deviceIDs = append(deviceIDs, d.ID)
	}
	defer func() {
		for _, id := range deviceIDs {
			repo.DeleteDevice(ctx, id)
		}
	}()

	fw, err := r.CreateFirmware(ctx, &model.CreateFirmwareRequest{Version: "ITEST-T1781408501-batch-v1.0.0"})
	if err != nil {
		t.Fatalf("CreateFirmware failed: %v", err)
	}
	defer r.DeleteFirmware(ctx, fw.ID)

	updates, err := r.BatchCreateUpdates(ctx, fw.ID, "batch-campaign", "", fw.Version, deviceIDs, 1, 3)
	if err != nil {
		t.Fatalf("BatchCreateUpdates failed: %v", err)
	}
	if len(updates) != 3 {
		t.Errorf("expected 3 updates, got %d", len(updates))
	}
	for _, u := range updates {
		if u.CampaignNameStr() != "batch-campaign" {
			t.Errorf("campaign = %s, want batch-campaign", u.CampaignNameStr())
		}
	}
}
