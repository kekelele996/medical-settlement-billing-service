package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func expectNoPanic(t *testing.T, name string, fn func() error) error {
	t.Helper()
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("%s panicked: %v", name, r)
			}
		}()
		err = fn()
	}()
	return err
}

func TestReverseSettlementChain(t *testing.T) {
	db := newTestDB(t)
	clientID, personID := seedData(t, db)
	ctx := context.Background()

	now := time.Now()
	db.Create(&model.UploadBatch{BatchNo: "B1", ClientID: clientID, InsuredPersonID: personID, UploadStatus: constants.UploadValidated, TotalAmount: 1000, ItemCount: 2})
	db.Create(&model.Presettlement{BatchID: 1, InsuredPersonID: personID, TotalAmount: 1000, InsurancePayAmount: 368})
	db.Create(&model.SettlementOrder{SettlementNo: "S1", BatchID: 1, InsuredPersonID: personID, PresettlementID: 1, ClientID: clientID, Status: constants.SettlementSettled, TotalAmount: 1000, InsurancePayAmount: 368, SettledAt: &now})

	insurance := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	batchRepo := repository.NewUploadBatchRepository(db)
	feeRepo := repository.NewFeeItemRepository(db)
	presetRepo := repository.NewPresettlementRepository(db)
	orderRepo := repository.NewSettlementOrderRepository(db)
	svc := NewSettlementService(presetRepo, orderRepo, feeRepo, batchRepo, insurance, util.NewSettlementCalculator(), testLogger())

	reversed, err := svc.ReverseSettlement(ctx, "S1")
	if err != nil {
		t.Fatalf("ReverseSettlement: %v", err)
	}
	if reversed.Status != constants.SettlementReversed {
		t.Fatalf("status = %s, want reversed", reversed.Status)
	}
	if _, err := svc.ReverseSettlement(ctx, "S1"); err == nil {
		t.Fatal("double reverse should fail")
	}

	err = expectNoPanic(t, "ReverseSettlement missing", func() error {
		_, err := svc.ReverseSettlement(ctx, "MISSING")
		return err
	})
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("ReverseSettlement missing = %v, want ErrNotFound", err)
	}
}
