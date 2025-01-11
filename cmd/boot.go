package cmd

import (
	"awesomeProject2/rebate/app/rebate/delivery/http"
	"awesomeProject2/rebate/app/rebate/repository"
	"awesomeProject2/rebate/app/rebate/usecase"
	"awesomeProject2/rebate/domain"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

// TODO: Fix database
var rebateUsecase domain.RebateUsecase
var dynamoClient *dynamodb.Client

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(repository.RebateProgram{}, repository.Transaction{}, repository.RebateClaim{})
	if err != nil {
		log.Fatal(err)
	}

}

func boot(db *dynamodb.Client) {
	var err error
	rebateRepository, err := repository.New(db)
	if err != nil {
		panic(err)
	}
	rebateUsecase, err = usecase.New(rebateRepository)
	if err != nil {
		panic(err)
	}

}

func Boot() {

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1")) // Update region if necessary
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(awsCfg)
	fmt.Println("Connected to DynamoDB.")
	ctx := context.TODO()
	createTables(ctx, dynamoClient)
	// Initialize DynamoDB client

	//migrate(dynamoClient)

	boot(dynamoClient)
	router := gin.Default()
	http.New(router.Group("api/"), rebateUsecase)
	router.StaticFile("/", "frontend/index.html")

	fmt.Println("project is up")
	// Start WebSocket server for real-time updates

	router.Run(":8080")

}
