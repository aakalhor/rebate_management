package repository

import (
	"awesomeProject2/rebate/domain"
	"github.com/google/uuid"
	"time"
)

type RebateProgram struct {
	ID                  string    `dynamodbav:"ID"`
	ProgramName         string    `dynamodbav:"ProgramName"`
	Percentage          float64   `dynamodbav:"Percentage"`
	StartDate           time.Time `dynamodbav:"StartDate"`
	EndDate             time.Time `dynamodbav:"EndDate"`
	EligibilityCriteria bool      `dynamodbav:"EligibilityCriteria"`
}

type RebateClaim struct {
	ID            string    `dynamodbav:"ID"`
	Amount        float64   `dynamodbav:"Amount"`
	TransactionID string    `dynamodbav:"TransactionID"`
	Status        string    `dynamodbav:"Status"`
	Date          time.Time `dynamodbav:"Date"`
}

type Transaction struct {
	ID       string    `dynamodbav:"ID"`
	Amount   float64   `dynamodbav:"Amount"`
	Date     time.Time `dynamodbav:"Date"`
	RebateID string    `dynamodbav:"RebateID"`
}

type CachedReport struct {
	CacheKey       string `dynamodbav:"CacheKey"`
	ReportData     string `dynamodbav:"ReportData"`
	ExpirationTime int64  `dynamodbav:"ExpirationTime"` // TTL attribute
}

func (transaction *Transaction) toDomain() *domain.Transaction {
	return &domain.Transaction{
		ID:       uuid.MustParse(transaction.ID),
		Amount:   transaction.Amount,
		Date:     transaction.Date,
		RebateID: uuid.MustParse(transaction.RebateID),
	}
}

func (rebateProgram *RebateProgram) toDomain() *domain.RebateProgram {
	return &domain.RebateProgram{
		ID:                  uuid.MustParse(rebateProgram.ID),
		ProgramName:         rebateProgram.ProgramName,
		Percentage:          rebateProgram.Percentage,
		StartDate:           rebateProgram.StartDate,
		EndDate:             rebateProgram.EndDate,
		EligibilityCriteria: rebateProgram.EligibilityCriteria,
	}
}

func (rebateClaim *RebateClaim) toDomain() *domain.RebateClaim {
	return &domain.RebateClaim{
		ID:            uuid.MustParse(rebateClaim.ID),
		Amount:        rebateClaim.Amount,
		TransactionID: uuid.MustParse(rebateClaim.TransactionID),
		Status:        domain.ClaimStatus(rebateClaim.Status),
		Date:          rebateClaim.Date,
	}
}
