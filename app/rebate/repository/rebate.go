package repository

import (
	"awesomeProject2/rebate/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"time"
)

type rebateRepository struct {
	db *dynamodb.Client
}

func (r *rebateRepository) ModifyClaimStatus(ctx context.Context, claimId uuid.UUID, status domain.ClaimStatus) (*domain.RebateClaim, error) {
	// Prepare the update expression and attribute values
	updateExpression := "SET #status = :status"
	expressionAttributeNames := map[string]string{
		"#status": "Status",
	}
	expressionAttributeValues := map[string]types.AttributeValue{
		":status": &types.AttributeValueMemberS{Value: string(status)}, // Convert domain.ClaimStatus to string
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("RebateClaim"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: claimId.String()}, // Convert uuid.UUID to string
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	// Perform the update
	result, err := r.db.UpdateItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to modify claim status: %w", err)
	}

	// Unmarshal the updated item back to domain.RebateClaim
	var updatedClaim domain.RebateClaim
	err = attributevalue.UnmarshalMap(result.Attributes, &updatedClaim)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated claim: %w", err)
	}

	return &updatedClaim, nil
}

func (r *rebateRepository) GetClaimByTransactionId(ctx context.Context, id uuid.UUID) (*domain.RebateClaim, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("RebateClaim"),
		FilterExpression: aws.String("TransactionID = :transactionID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":transactionID": &types.AttributeValueMemberS{Value: id.String()},
		},
	}

	result, err := r.db.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rebate claim by transaction ID: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, domain.ErrClaimNotFound
	}

	var repoClaim RebateClaim
	err = attributevalue.UnmarshalMap(result.Items[0], &repoClaim)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal rebate claim: %w", err)
	}

	return repoClaim.toDomain(), nil
}

func (r *rebateRepository) StoreRebateProgram(ctx context.Context, program domain.RebateProgram) (*domain.RebateProgram, error) {
	// Check if ProgramName is unique
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("RebateProgram"),
		IndexName:              aws.String("ProgramNameUniqueIndex"), // Use GSI for ProgramName
		KeyConditionExpression: aws.String("ProgramName = :programName"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":programName": &types.AttributeValueMemberS{Value: program.ProgramName},
		},
	}

	queryResult, err := r.db.Query(ctx, queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query RebateProgram table for uniqueness: %w", err)
	}

	if len(queryResult.Items) > 0 {
		return nil, fmt.Errorf("program name '%s' already exists", program.ProgramName)
	}

	// Convert domain model to repository model
	repoProgram := RebateProgram{
		ID:                  program.ID.String(),
		ProgramName:         program.ProgramName,
		Percentage:          program.Percentage,
		StartDate:           program.StartDate,
		EndDate:             program.EndDate,
		EligibilityCriteria: program.EligibilityCriteria,
	}

	// Marshal to DynamoDB item
	item, err := attributevalue.MarshalMap(repoProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rebate program: %w", err)
	}

	// Write to DynamoDB
	input := &dynamodb.PutItemInput{
		TableName: aws.String("RebateProgram"),
		Item:      item,
	}

	_, err = r.db.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to store rebate program: %w", err)
	}

	return &program, nil
}

func (r *rebateRepository) StoreTransaction(ctx context.Context, transaction domain.Transaction) (*domain.Transaction, error) {
	tempTransaction, err := r.GetTransactionByID(ctx, transaction.ID)
	if err != nil {
		if err != domain.ErrTransactionNotFound {
			return nil, err
		}
	}

	if tempTransaction != nil {
		return nil, domain.ErrTransactionCanNotCreate
	}
	repoTransaction := Transaction{
		ID:       transaction.ID.String(), // Convert uuid.UUID to string
		Amount:   transaction.Amount,
		Date:     transaction.Date,
		RebateID: transaction.RebateID.String(), // Convert uuid.UUID to string
	}

	item, err := attributevalue.MarshalMap(repoTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Transaction"),
		Item:      item,
	}

	_, err = r.db.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to store transaction: %w", err)
	}

	return &transaction, nil
}

func (r *rebateRepository) GetRebateByID(ctx context.Context, id uuid.UUID) (*domain.RebateProgram, error) {
	// DynamoDB query
	input := &dynamodb.GetItemInput{
		TableName: aws.String("RebateProgram"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id.String()}, // Use uuid.String() for querying
		},
	}

	result, err := r.db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, domain.ErrRebateNotFound
	}

	// Unmarshal to repository model
	var repoProgram RebateProgram
	err = attributevalue.UnmarshalMap(result.Item, &repoProgram)
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	return repoProgram.toDomain(), nil
}

func (r *rebateRepository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Transaction"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id.String()},
		},
	}

	result, err := r.db.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	if result.Item == nil {
		return nil, domain.ErrTransactionNotFound
	}

	var transaction Transaction

	err = attributevalue.UnmarshalMap(result.Item, &transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return transaction.toDomain(), nil
}

func (r *rebateRepository) StoreRebateClaim(ctx context.Context, claim domain.RebateClaim) (*domain.RebateClaim, error) {
	repoClaim := RebateClaim{
		ID:            claim.ID.String(), // Convert uuid.UUID to string
		Amount:        claim.Amount,
		TransactionID: claim.TransactionID.String(), // Convert uuid.UUID to string
		Status:        string(claim.Status),
		Date:          claim.Date,
	}

	item, err := attributevalue.MarshalMap(repoClaim)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rebate claim: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("RebateClaim"),
		Item:      item,
	}

	_, err = r.db.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to store rebate claim: %w", err)
	}

	return &claim, nil
}

func (r *rebateRepository) ListClaimsWithinInterval(ctx context.Context, from time.Time, to time.Time) ([]domain.RebateClaim, error) {

	fromStr := from.Format(time.RFC3339)
	toStr := to.Format(time.RFC3339)

	// Prepare the QueryInput
	input := &dynamodb.ScanInput{
		TableName:        aws.String("RebateClaim"),
		FilterExpression: aws.String("#date BETWEEN :from AND :to"),
		ExpressionAttributeNames: map[string]string{
			"#date": "Date", // Use an alias for the Date attribute
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":from": &types.AttributeValueMemberS{Value: fromStr},
			":to":   &types.AttributeValueMemberS{Value: toStr},
		},
	}

	// Execute the scan
	result, err := r.db.Scan(ctx, input)
	if err != nil {
		return nil, domain.ErrFailedToListClaims
	}
	fmt.Println("pass here!!", result, err)
	// Check if there are no items
	if len(result.Items) == 0 {
		return nil, nil
	}

	// Unmarshal the result into a slice of RebateClaim
	var repoClaims []RebateClaim
	err = attributevalue.UnmarshalListOfMaps(result.Items, &repoClaims)
	if err != nil {
		return nil, err
	}

	// Convert to domain models
	domainClaims := make([]domain.RebateClaim, len(repoClaims))
	for i, repoClaim := range repoClaims {

		domainClaims[i] = *repoClaim.toDomain()
	}

	return domainClaims, nil
}

func (r *rebateRepository) GetCachedReport(ctx context.Context, cacheKey string) (*domain.RebateClaimsReport, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("ReportCache"),
		Key: map[string]types.AttributeValue{
			"CacheKey": &types.AttributeValueMemberS{Value: cacheKey},
		},
	}

	result, err := r.db.GetItem(ctx, input)
	if err != nil {
		return nil, domain.ErrFailedToGetCache
	}

	if result.Item == nil {
		return nil, nil // Cache miss
	}

	var cached CachedReport
	err = attributevalue.UnmarshalMap(result.Item, &cached)
	if err != nil {
		return nil, err
	}

	var report domain.RebateClaimsReport
	err = json.Unmarshal([]byte(cached.ReportData), &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *rebateRepository) StoreCachedReport(ctx context.Context, cacheKey string, report *domain.RebateClaimsReport, ttl time.Duration) error {
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return err
	}

	expirationTime := time.Now().Add(ttl).Unix()

	item := CachedReport{
		CacheKey:       cacheKey,
		ReportData:     string(reportJSON),
		ExpirationTime: expirationTime,
	}

	dynamoItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("ReportCache"),
		Item:      dynamoItem,
	}

	_, err = r.db.PutItem(ctx, input)
	if err != nil {
		return domain.ErrFailedToStoreCache
	}

	return nil
}

func New(db *dynamodb.Client) (*rebateRepository, error) {
	return &rebateRepository{
		db: db,
	}, nil
}

var _ domain.RebateRepository = &rebateRepository{}
