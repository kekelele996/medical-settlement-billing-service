package service

import (
	"context"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/repository"
)

func TestFeeUploadDuplicateAndSummary(t *testing.T) {
	db := newTestDB(t)
	clientID, personID := seedData(t, db)
	ctx := context.Background()

	insurance := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	batchRepo := repository.NewUploadBatchRepository(db)
	feeRepo := repository.NewFeeItemRepository(db)
	svc := NewFeeService(batchRepo, feeRepo, insurance, testLogger())

	items := []FeeItemInput{
		{ItemCode: "A1", ItemName: "甲类药品", ItemType: constants.FeeItemDrug, UnitPrice: 500, Quantity: 1, MedicalCategory: constants.MedicalCategoryClassA},
		{ItemCode: "A2", ItemName: "甲类检查", ItemType: constants.FeeItemExam, UnitPrice: 400, Quantity: 1, MedicalCategory: constants.MedicalCategoryClassA},
	}
	up, err := svc.Upload(ctx, UploadInput{ClientID: clientID, InsuredPersonID: personID, Items: items})
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if up.UploadStatus != constants.UploadValidated {
		t.Fatalf("upload status = %s, want validated", up.UploadStatus)
	}
	if up.ItemCount != 2 {
		t.Fatalf("ItemCount = %d, want 2", up.ItemCount)
	}
	if up.TotalAmount != 900 {
		t.Fatalf("TotalAmount = %v, want 900", up.TotalAmount)
	}

	if _, err := svc.Upload(ctx, UploadInput{ClientID: clientID, InsuredPersonID: personID, Items: items}); err == nil {
		t.Fatal("second upload same day should be rejected as duplicate")
	}

	persisted, err := feeRepo.ListByBatch(up.BatchID)
	if err != nil {
		t.Fatalf("ListByBatch: %v", err)
	}
	if len(persisted) != 2 || persisted[0].ItemCode != "A1" {
		t.Fatalf("items order/content wrong: %+v", persisted)
	}
}
