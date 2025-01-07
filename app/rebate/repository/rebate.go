package repository

import (
	"awesomeProject2/rebate/domain"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type rebateRepository struct {
	//	TODO: Fix database
	db *gorm.DB
}

func (r *rebateRepository) StoreRebateProgram(ctx context.Context, program domain.RebateProgram) (*domain.RebateProgram, error) {
	//TODO implement me
	return &domain.RebateProgram{}, nil
}

func (r *rebateRepository) StoreTransaction(ctx context.Context, program domain.Transaction) (*domain.Transaction, error) {
	//TODO implement me
	return &domain.Transaction{}, nil
}

func (r *rebateRepository) GetRebateByID(ctx context.Context, id uuid.UUID) (*domain.RebateProgram, error) {
	//TODO implement me
	return &domain.RebateProgram{}, nil
}

func (r *rebateRepository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	//TODO implement me
	return &domain.Transaction{}, nil
}

func (r *rebateRepository) StoreRebateClaim(ctx context.Context, transactionId uuid.UUID) (uuid.UUID, error) {
	//TODO implement me
	return uuid.New(), nil
}

func (r *rebateRepository) ListClaimsWithinInterval(ctx context.Context, from time.Time, to time.Time) ([]domain.RebateClaim, error) {
	//TODO implement me
	return []domain.RebateClaim{}, nil
}

func New(db *gorm.DB) (*rebateRepository, error) {
	return &rebateRepository{
		db: db,
	}, nil
}

var _ domain.RebateRepository = &rebateRepository{}
