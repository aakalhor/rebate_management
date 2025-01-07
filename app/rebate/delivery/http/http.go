package http

import (
	"awesomeProject2/rebate/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	rg.GET("/reporting", handler.reporting)

}

func (h *handler) createClaim(context *gin.Context) {
	transactionId := uuid.MustParse(context.GetString("transactionId"))
	claimId, err := h.u.SubmitRebateClaim(context, transactionId)
	if err != nil {
		//TODO: Fix errMap
		ErrHandler(context, err, errMap)
	}
	context.JSON(http.StatusCreated, &claimId)
}

func (h *handler) createRebate(context *gin.Context) {

	rebate := domain.RebateProgram{
		ID:                  uuid.New(),
		ProgramName:         "",
		Percentage:          0,
		StartDate:           time.Time{},
		EndDate:             time.Time{},
		EligibilityCriteria: false,
	}
	retRebate, err := h.u.CreateRebateProgram(context, rebate)
	if err != nil {
		//TODO: Fix errMap
		ErrHandler(context, err, errMap)
	}
	context.JSON(http.StatusCreated, &retRebate)
}

func (h *handler) createTransaction(context *gin.Context) {
	tu := domain.Transaction{
		ID:       uuid.New(),
		Amount:   0,
		Date:     time.Time{},
		RebateID: 0,
	}
	retTransaction, err := h.u.SubmitTransaction(context, tu)
	if err != nil {
		x.ErrHandler(context, err, errMap)
	}

	context.JSON(http.StatusCreated, &retTransaction)
}

func (h *handler) calculateRebate(context *gin.Context) {

}

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
		x.ErrHandler(context, err, errMap)
	}

	context.JSON(http.StatusCreated, &rci)
}
