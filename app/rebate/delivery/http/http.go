package http

import (
	"awesomeProject2/rebate/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

type handler struct {
	u domain.RebateUsecase
}

func New(rg *gin.RouterGroup, u domain.RebateUsecase) {
	handler := &handler{
		u: u,
	}
	rg.POST("/rebate", handler.createRebate)
	rg.POST("/transaction", handler.createTransaction)
	rg.POST("/claim", handler.createClaim)
	rg.GET("/calculate", handler.calculateRebate)
	rg.POST("/reporting", handler.reporting)
	rg.GET("/claims/status", handler.getClaimStatus)

	rg.PUT("/claim", handler.changeClaimStatus)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (h *handler) getClaimStatus(context *gin.Context) {
	// Call the use case to get claim status report
	report, err := h.u.ReportClaimsByPeriod(context, time.Time{}, time.Now())
	if err != nil {
		ErrHandler(context, err, errMap)
		return
	}

	// Return the aggregated status counts
	context.JSON(http.StatusOK, gin.H{
		"total":    report.Total.Count,
		"pending":  report.Pending.Count,
		"approved": report.Approved.Count,
		"rejected": report.Rejected.Count,
	})
}

// Create Claim godoc
// @Summary Submit a rebate claim
// @Schemes http https
// @Description Submit a claim using a valid transaction ID. The claim must meet eligibility criteria.
// @Tags Rebate Management
// @Accept json
// @Produce json
// @Param transaction_id query string true "Transaction ID (UUID format)"
// @Success 201 {object} RebateClaim "Claim successfully created"
// @Failure 400 {object} CodeResponse "Invalid input or user not eligible"
// @Failure 404 {object} CodeResponse "Transaction not found"
// @Failure 500 {object} CodeResponse "Internal server error"
// @Router /api/claim [post]
func (h *handler) createClaim(context *gin.Context) {
	transactionId, err := uuid.Parse(context.Query("transaction_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	claimId := uuid.New()
	date := time.Now()
	claim, err := h.u.SubmitRebateClaim(context, claimId, transactionId, date)
	if err != nil {
		ErrHandler(context, err, errMap)
		return
	}
	context.JSON(http.StatusCreated, RebateClaimFromDomain(*claim))
}

// Change Claim godoc
// @Summary Change rebate claim status
// @Schemes http https
// @Description Submit a claim using a valid claim ID. The claim must meet eligibility criteria.
// @Tags Rebate Management
// @Accept json
// @Produce json
// @Param claim_id query string true "Claim ID (UUID format)"
// @Success 202 {object} RebateClaim "Status successfully changed"
// @Failure 400 {object} CodeResponse "Invalid input for uuid or status"
// @Failure 500 {object} CodeResponse "Internal server error"
// @Router /api/claim [put]
func (h *handler) changeClaimStatus(context *gin.Context) {
	claimId, err := uuid.Parse(context.Query("claim_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid uuid"})
		return
	}

	status := context.Query("status")
	if status == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "status field is needed"})
		return
	}

	// Validate if status is a valid ClaimStatus
	validStatuses := map[domain.ClaimStatus]bool{
		domain.StatusPending:  true,
		domain.StatusApproved: true,
		domain.StatusRejected: true,
	}
	claimStatus := domain.ClaimStatus(status)
	if !validStatuses[claimStatus] {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid status value"})
		return
	}
	_, err = h.u.ChangeClaimStatus(context, claimId, claimStatus)
	if err != nil {
		ErrHandler(context, err, errMap)
		return
	}
	context.JSON(http.StatusAccepted, gin.H{})
}

// Create Rebate Program godoc
// @Summary Register a new rebate program
// @Schemes http https
// @Description Register a rebate program. Program names must be unique and include a percentage rebate.
// @Tags Rebate Management
// @Accept json
// @Produce json
// @Param rebate body RebateProgram true "Rebate Program Details"
// @Success 201 {object} RebateProgram "Rebate program successfully created"
// @Failure 400 {object} CodeResponse "Invalid input or duplicate program name"
// @Failure 500 {object} CodeResponse "Internal server error"
// @Router /api/rebate [post]
func (h *handler) createRebate(context *gin.Context) {
	var r RebateProgram

	if err := context.BindJSON(&r); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if r.Percentage > 100 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid rebate percentage"})
		return
	}

	rebate := domain.RebateProgram{
		ID:                  uuid.New(),
		ProgramName:         r.ProgramName,
		Percentage:          r.Percentage,
		StartDate:           r.StartDate,
		EndDate:             r.EndDate,
		EligibilityCriteria: r.EligibilityCriteria,
	}

	retRebate, err := h.u.CreateRebateProgram(context, rebate)
	if err != nil {
		ErrHandler(context, err, errMap)
		return
	}
	context.JSON(http.StatusCreated, RebateProgramFromDomain(*retRebate))
}

// Create Transaction godoc
// @Summary Record a new transaction
// @Schemes http https
// @Description Record a transaction for a rebate program. The transaction must reference an existing rebate program.
// @Tags Transaction Management
// @Accept json
// @Produce json
// @Param transaction body Transaction true "Transaction Details"
// @Success 201 {object} Transaction "Transaction successfully recorded"
// @Failure 400 {object} CodeResponse "Invalid input or unable to create transaction"
// @Failure 404 {object} CodeResponse "Rebate program not found"
// @Failure 500 {object} CodeResponse "Internal server error"
// @Router /api/transaction [post]
func (h *handler) createTransaction(context *gin.Context) {
	var t Transaction

	if err := context.BindJSON(&t); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	transaction := domain.Transaction{
		ID:       uuid.New(),
		Amount:   t.Amount,
		Date:     t.Date,
		RebateID: t.RebateID,
	}

	retTransaction, err := h.u.SubmitTransaction(context, transaction)
	if err != nil {
		ErrHandler(context, err, errMap)
		return
	}

	context.JSON(http.StatusCreated, TransactionFromDomain(*retTransaction))
}

// Calculate Rebate godoc
// @Summary Calculate rebate amount
// @Schemes http https
// @Description Calculate the rebate amount for a transaction. Requires a valid transaction ID.
// @Tags Rebate Management
// @Accept json
// @Produce json
// @Param transaction_id query string true "Transaction ID (UUID format)"
// @Success 200 {object} map[string]float64 "Rebate amount successfully calculated"
// @Failure 400 {object} CodeResponse "Invalid transaction or rebate calculation error"
// @Failure 404 {object} CodeResponse "Transaction not found"
// @Failure 500 {object} CodeResponse "Internal server error"
// @Router /api/calculate [get]
func (h *handler) calculateRebate(context *gin.Context) {
	transactionId, err := uuid.Parse(context.Query("transaction_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	rebateAmount, err := h.u.CalculateRebateOfTransaction(context, transactionId)
	if err != nil {
		ErrHandler(context, err, errMap)
		return
	}

	context.JSON(http.StatusOK, gin.H{"rebate_amount": rebateAmount})
}

// Report Claims godoc
// @Summary Generate claims report
// @Schemes http https
// @Description Generate a detailed report of claims within a specified date range.
// @Tags Reporting
// @Accept json
// @Produce json
// @Param dateRange body DateRange true "Date Range (start_date and end_date in YYYY-MM-DD format)"
// @Success 200 {object} RebateClaimsReport "Claims report successfully generated"
// @Failure 400 {object} CodeResponse "Invalid date range or input"
// @Failure 500 {object} CodeResponse "Internal server error"
// @Router /api/reporting [post]
func (h *handler) reporting(context *gin.Context) {

	var dateRange DateRange
	if err := context.ShouldBindJSON(&dateRange); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO: Fix errors on dates
	startDate, err := time.Parse("2006-01-02", dateRange.StartDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", dateRange.EndDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	// Ensure start_date is earlier than end_date
	if !startDate.Before(endDate) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be earlier than end_date"})
		return
	}
	rci, err := h.u.ReportClaimsByPeriod(context, startDate, endDate)
	if err != nil {
		ErrHandler(context, err, errMap)
		return
	}

	context.JSON(http.StatusOK, RebateClaimReportFromDomain(*rci))
}
