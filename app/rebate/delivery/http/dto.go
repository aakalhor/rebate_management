package http

import (
	"awesomeProject2/rebate/domain"
	"github.com/google/uuid"
	"time"
)

type RebateProgram struct {
	ID                  uuid.UUID `json:"id"`
	ProgramName         string    `json:"program_name" binding:"required"`
	Percentage          float64   `json:"percentage" binding:"required"`
	StartDate           time.Time `json:"start_date" binding:"required"`
	EndDate             time.Time `json:"end_date" binding:"required"`
	EligibilityCriteria bool      `json:"eligibility_criteria" binding:"required"`
}

type Transaction struct {
	ID       uuid.UUID `json:"id"`
	Amount   float64   `json:"amount"`
	Date     time.Time `json:"date"`
	RebateID uuid.UUID `json:"rebate_id"`
}

type RebateClaim struct {
	ID            uuid.UUID `json:"id"`
	Amount        float64   `json:"amount"`
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	Date          time.Time `json:"date"`
}

type RebateClaimsReport struct {
	From     time.Time    `json:"from"`
	To       time.Time    `json:"to"`
	Total    ClaimMetrics `json:"total"`
	Pending  ClaimMetrics `json:"pending"`
	Approved ClaimMetrics `json:"approved"`
	Rejected ClaimMetrics `json:"rejected"`
}

type ClaimMetrics struct {
	Count  uint64  `json:"count"`
	Amount float64 `json:"amount"`
}

type DateRange struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

func RebateProgramFromDomain(rebateProgram domain.RebateProgram) RebateProgram {

	return RebateProgram{
		ID:                  rebateProgram.ID,
		ProgramName:         rebateProgram.ProgramName,
		Percentage:          rebateProgram.Percentage,
		StartDate:           rebateProgram.StartDate,
		EndDate:             rebateProgram.EndDate,
		EligibilityCriteria: rebateProgram.EligibilityCriteria,
	}
}

func TransactionFromDomain(transaction domain.Transaction) Transaction {

	return Transaction{
		ID:       transaction.ID,
		Amount:   transaction.Amount,
		Date:     transaction.Date,
		RebateID: transaction.RebateID,
	}
}

func RebateClaimFromDomain(rebateClaim domain.RebateClaim) RebateClaim {

	return RebateClaim{
		ID:            rebateClaim.ID,
		Amount:        rebateClaim.Amount,
		TransactionID: rebateClaim.TransactionID,
		Status:        string(rebateClaim.Status),
		Date:          rebateClaim.Date,
	}
}

func RebateClaimReportFromDomain(rebateClaimsReport domain.RebateClaimsReport) RebateClaimsReport {

	return RebateClaimsReport{
		From:     rebateClaimsReport.From,
		To:       rebateClaimsReport.To,
		Total:    ClaimMetrics(rebateClaimsReport.Total),
		Pending:  ClaimMetrics(rebateClaimsReport.Pending),
		Approved: ClaimMetrics(rebateClaimsReport.Approved),
		Rejected: ClaimMetrics(rebateClaimsReport.Rejected),
	}
}
