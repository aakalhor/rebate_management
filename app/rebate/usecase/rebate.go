package usecase

import (
	"awesomeProject2/rebate/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type rebateUsecase struct {
	r domain.RebateRepository
}

func (r *rebateUsecase) ChangeClaimStatus(ctx context.Context, claimId uuid.UUID, status domain.ClaimStatus) (*domain.RebateClaim, error) {
	return r.r.ModifyClaimStatus(ctx, claimId, status)
}

func (r *rebateUsecase) CreateRebateProgram(ctx context.Context, program domain.RebateProgram) (*domain.RebateProgram, error) {
	return r.r.StoreRebateProgram(ctx, program)
}

func (r *rebateUsecase) SubmitTransaction(ctx context.Context, transaction domain.Transaction) (*domain.Transaction, error) {

	_, err := r.r.GetRebateByID(ctx, transaction.RebateID)
	if err != nil {
		return nil, err
	}
	return r.r.StoreTransaction(ctx, transaction)
}

func (r *rebateUsecase) CalculateRebateOfTransaction(ctx context.Context, transactionId uuid.UUID) (float64, error) {

	transaction, err := r.r.GetTransactionByID(ctx, transactionId)
	if err != nil {
		return 0, err
	}
	if transaction.ID == uuid.Nil {
		return 0, domain.ErrTransactionNotFound
	}
	rebate, err := r.r.GetRebateByID(ctx, transaction.RebateID)
	if err != nil {
		return 0, err
	}
	if transaction.ID == uuid.Nil {
		return 0, domain.ErrRebateNotFound
	}

	if rebate.EligibilityCriteria == false {
		return 0, domain.ErrNotEligible
	}

	if transaction.Date.Before(rebate.StartDate) || transaction.Date.After(rebate.EndDate) {
		return 0, domain.ErrInvalidInterval
	}

	rebateAmount := rebate.Percentage / 100 * transaction.Amount
	return rebateAmount, nil
}

func (r *rebateUsecase) ReportClaimsByPeriod(ctx context.Context, from time.Time, to time.Time) (*domain.RebateClaimsReport, error) {
	calims, err := r.r.ListClaimsWithinInterval(ctx, from, to)
	if err != nil {
		return nil, err
	}
	var totalClaimCount = len(calims)
	var (
		pendingClaimCount   uint64
		pendingClaimAmount  float64
		approvedClaimCount  uint64
		approvedClaimAmount float64
		rejectedClaimCount  uint64
		rejectedClaimAmount float64
	)

	for _, claim := range calims {
		switch claim.Status {
		case domain.StatusPending:
			pendingClaimCount += 1
			pendingClaimAmount += claim.Amount
		case domain.StatusApproved:
			approvedClaimCount += 1
			approvedClaimAmount += claim.Amount
		case domain.StatusRejected:
			rejectedClaimCount += 1
			rejectedClaimAmount += claim.Amount
		default:
			return &domain.RebateClaimsReport{}, fmt.Errorf("unexpected claim status: %s", claim.Status)
		}
	}
	report := domain.RebateClaimsReport{
		From: from,
		To:   to,
		Total: domain.ClaimMetrics{
			Count:  uint64(totalClaimCount),
			Amount: pendingClaimAmount + approvedClaimAmount + rejectedClaimAmount,
		},
		Pending: domain.ClaimMetrics{
			Count:  pendingClaimCount,
			Amount: pendingClaimAmount,
		},
		Approved: domain.ClaimMetrics{
			Count:  approvedClaimCount,
			Amount: approvedClaimAmount,
		},
		Rejected: domain.ClaimMetrics{
			Count:  rejectedClaimCount,
			Amount: rejectedClaimAmount,
		},
	}

	return &report, nil

}

func (r *rebateUsecase) SubmitRebateClaim(ctx context.Context, claimId uuid.UUID, transactionId uuid.UUID, date time.Time) (*domain.RebateClaim, error) {
	tempClaim, err := r.r.GetClaimByTransactionId(ctx, transactionId)

	if err != nil && err != domain.ErrClaimNotFound {
		return &domain.RebateClaim{}, err
	}
	fmt.Println(tempClaim, err, "LLLAS")
	if err != domain.ErrClaimNotFound {
		return tempClaim, domain.ErrRebateAlreadyClaimed
	}

	transaction, err := r.r.GetTransactionByID(ctx, transactionId)
	if err != nil {
		return nil, err
	}
	if transaction.ID == uuid.Nil {
		return nil, domain.ErrTransactionNotFound
	}

	claim := domain.RebateClaim{
		ID:            claimId,
		Amount:        transaction.Amount,
		TransactionID: transactionId,
		Status:        domain.StatusPending,
		Date:          date,
	}

	return r.r.StoreRebateClaim(ctx, claim)
}

func New(r domain.RebateRepository) (*rebateUsecase, error) {
	return &rebateUsecase{
		r: r,
	}, nil
}

var _ domain.RebateUsecase = &rebateUsecase{}
