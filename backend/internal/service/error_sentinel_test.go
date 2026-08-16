package service

import (
	"context"
	"errors"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func TestNotFoundErrorsRemainNotFound(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	insurance := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	presetRepo := repository.NewPresettlementRepository(db)
	orderRepo := repository.NewSettlementOrderRepository(db)
	feeRepo := repository.NewFeeItemRepository(db)
	batchRepo := repository.NewUploadBatchRepository(db)
	svc := NewSettlementService(presetRepo, orderRepo, feeRepo, batchRepo, insurance, util.NewSettlementCalculator(), testLogger())

	_, err := insurance.GetByID(ctx, 99999)
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("GetByID missing = %v, want ErrNotFound", err)
	}
	if appErr, ok := err.(*util.AppError); !ok || appErr.Code != 1003 {
		t.Fatalf("GetByID missing code = %v, want 1003", err)
	}

	_, err = svc.GetOrder(ctx, "MISSING")
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("GetOrder missing = %v, want ErrNotFound", err)
	}

	_, err = svc.ReverseSettlement(ctx, "MISSING")
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("ReverseSettlement missing = %v, want ErrNotFound", err)
	}
}
