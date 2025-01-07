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
	rg.POST("/claim", createClaim)
	rg.GET("/calculate", handler.calculateRebate)
	rg.GET("/reporting", handler.reporting)

}

func createClaim(context *gin.Context) {

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

}
