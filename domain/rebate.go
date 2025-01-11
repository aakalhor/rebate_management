package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID       uuid.UUID
	Amount   float64
	Date     time.Time
	RebateID uuid.UUID
}

type RebateProgram struct {
	ID                  uuid.UUID
	ProgramName         string
	Percentage          float64
	StartDate           time.Time
	EndDate             time.Time
	EligibilityCriteria bool
}

type ClaimStatus string

const (
	StatusPending  ClaimStatus = "pending"
	StatusApproved ClaimStatus = "approved"
	StatusRejected ClaimStatus = "rejected"
)

type RebateClaim struct {
	ID            uuid.UUID
	Amount        float64
	TransactionID uuid.UUID
	Status        ClaimStatus
	Date          time.Time
}

type ClaimMetrics struct {
	Count  uint64
	Amount float64
}

type RebateClaimsReport struct {
	From     time.Time
	To       time.Time
	Total    ClaimMetrics
	Pending  ClaimMetrics
	Approved ClaimMetrics
	Rejected ClaimMetrics
}

type RebateUsecase interface {
	CreateRebateProgram(ctx context.Context, program RebateProgram) (*RebateProgram, error)
	SubmitTransaction(ctx context.Context, transaction Transaction) (*Transaction, error)
	CalculateRebateOfTransaction(ctx context.Context, transactionId uuid.UUID) (float64, error)
	ReportClaimsByPeriod(ctx context.Context, from time.Time, to time.Time) (*RebateClaimsReport, error)
	SubmitRebateClaim(ctx context.Context, claimId uuid.UUID, transactionId uuid.UUID, date time.Time) (*RebateClaim, error)
	ChangeClaimStatus(ctx context.Context, claimId uuid.UUID, status ClaimStatus) (*RebateClaim, error)
}

type RebateRepository interface {
	StoreRebateProgram(ctx context.Context, program RebateProgram) (*RebateProgram, error)
	StoreTransaction(ctx context.Context, program Transaction) (*Transaction, error)
	GetRebateByID(ctx context.Context, id uuid.UUID) (*RebateProgram, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetClaimByTransactionId(ctx context.Context, id uuid.UUID) (*RebateClaim, error)
	StoreRebateClaim(ctx context.Context, calim RebateClaim) (*RebateClaim, error)
	ListClaimsWithinInterval(ctx context.Context, from time.Time, to time.Time) ([]RebateClaim, error)
	GetCachedReport(ctx context.Context, cacheKey string) (*RebateClaimsReport, error)
	StoreCachedReport(ctx context.Context, cacheKey string, report *RebateClaimsReport, ttl time.Duration) error
	ModifyClaimStatus(ctx context.Context, claimId uuid.UUID, status ClaimStatus) (*RebateClaim, error)
}
