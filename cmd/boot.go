package cmd

import (
	"awesomeProject2/rebate/app/rebate/delivery/http"
	"awesomeProject2/rebate/app/rebate/repository"
	"awesomeProject2/rebate/app/rebate/usecase"
	"awesomeProject2/rebate/domain"
	"github.com/gin-gonic/gin"
)

// TODO: Fix database
var rebateUsecase domain.RebateUsecase

func boot(db *db.Database) {
	var err error
	rebateRepository, err := repository.New(db.GormDB)
	if err != nil {
		panic(err)
	}
	rebateUsecase, err = usecase.New(rebateRepository)
	if err != nil {
		panic(err)
	}

}

func Boot() {
	// *****run http server*****
	db := Database
	boot(db)
	router := gin.Default()
	http.New(router.Group("order/"), rebateUsecase)

	router.Run(":8080")
}
