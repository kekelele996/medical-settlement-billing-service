package service

import (
	"context"
	"errors"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func expectNoPanicErr(t *testing.T, name string, fn func() error) error {
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

func TestMissingPersonAndOrderReturnNotFound(t *testing.T) {
	db := newTestDB(t)
	clientID, _ := seedData(t, db)
	ctx := context.Background()

	db.Create(&model.UploadBatch{BatchNo: "B1", ClientID: clientID, InsuredPersonID: 99999, UploadStatus: constants.UploadValidated, TotalAmount: 1000, ItemCount: 2})

	insurance := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	batchRepo := repository.NewUploadBatchRepository(db)
	feeRepo := repository.NewFeeItemRepository(db)
	presetRepo := repository.NewPresettlementRepository(db)
	orderRepo := repository.NewSettlementOrderRepository(db)
	svc := NewSettlementService(presetRepo, orderRepo, feeRepo, batchRepo, insurance, util.NewSettlementCalculator(), testLogger())

	err := expectNoPanicErr(t, "CalculatePresettlement", func() error {
		_, err := svc.CalculatePresettlement(ctx, 1)
		return err
	})
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("CalculatePresettlement missing person = %v, want ErrNotFound", err)
	}

	err = expectNoPanicErr(t, "ReverseSettlement", func() error {
		_, err := svc.ReverseSettlement(ctx, "MISSING")
		return err
	})
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("ReverseSettlement missing = %v, want ErrNotFound", err)
	}
}
