package service

import (
	"context"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func TestPresettlementCalcChain(t *testing.T) {
	db := newTestDB(t)
	clientID, personID := seedData(t, db)
	ctx := context.Background()

	insurance := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	batchRepo := repository.NewUploadBatchRepository(db)
	feeRepo := repository.NewFeeItemRepository(db)
	feeSvc := NewFeeService(batchRepo, feeRepo, insurance, testLogger())

	up, err := feeSvc.Upload(ctx, UploadInput{
		ClientID:        clientID,
		InsuredPersonID: personID,
		Items: []FeeItemInput{
			{ItemCode: "A1", ItemName: "甲类药品", ItemType: constants.FeeItemDrug, UnitPrice: 600, Quantity: 1, MedicalCategory: constants.MedicalCategoryClassA},
			{ItemCode: "B1", ItemName: "乙类检查", ItemType: constants.FeeItemExam, UnitPrice: 400, Quantity: 1, MedicalCategory: constants.MedicalCategoryClassB},
		},
	})
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if up.UploadStatus != constants.UploadValidated {
		t.Fatalf("upload status = %s, want validated", up.UploadStatus)
	}

	presetRepo := repository.NewPresettlementRepository(db)
	orderRepo := repository.NewSettlementOrderRepository(db)
	calc := util.NewSettlementCalculator()
	svc := NewSettlementService(presetRepo, orderRepo, feeRepo, batchRepo, insurance, calc, testLogger())

	preset, err := svc.CalculatePresettlement(ctx, up.BatchID)
	if err != nil {
		t.Fatalf("CalculatePresettlement: %v", err)
	}
	if preset.TotalAmount != 1000 {
		t.Fatalf("TotalAmount = %v, want 1000", preset.TotalAmount)
	}
	if preset.InsurancePayAmount != 368 {
		t.Fatalf("InsurancePayAmount = %v, want 368", preset.InsurancePayAmount)
	}
	if preset.PersonalAccountAmount != 368 {
		t.Fatalf("PersonalAccountAmount = %v, want 368", preset.PersonalAccountAmount)
	}
	if preset.SelfPayAmount != 632 {
		t.Fatalf("SelfPayAmount = %v, want 632", preset.SelfPayAmount)
	}
}

