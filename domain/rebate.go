package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID       uuid.UUID `json:"primaryKey"`
	Amount   float64   `json:"unique;not null"`
	Date     time.Time `json:"not null"`
	RebateID uuid.UUID `json:"not null"`
}

type RebateProgram struct {
	ID                  uuid.UUID `gorm:"primaryKey"`
	ProgramName         string    `gorm:"unique;not null"`
	Percentage          float64   `gorm:"not null"`
	StartDate           time.Time `gorm:"not null"`
	EndDate             time.Time `gorm:"not null"`
	EligibilityCriteria bool      `gorm:"bool"`
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
	SubmitRebateClaim(ctx context.Context, transactionId uuid.UUID) (claimId uuid.UUID, err error)
}

type RebateRepository interface {
	StoreRebateProgram(ctx context.Context, program RebateProgram) (*RebateProgram, error)
	StoreTransaction(ctx context.Context, program Transaction) (*Transaction, error)
	GetRebateByID(ctx context.Context, id uuid.UUID) (*RebateProgram, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	StoreRebateClaim(ctx context.Context, transactionId uuid.UUID) (claimId uuid.UUID, err error)
	ListClaimsWithinInterval(ctx context.Context, from time.Time, to time.Time) ([]RebateClaim, error)
}
