package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func createTables(ctx context.Context, dynamoClient *dynamodb.Client) {
	tables := []struct {
		Name string
		Key  string
	}{
		{"RebateProgram", "ProgramName"},
		{"RebateClaim", "TransactionID"},
		{"Transaction", "ID"},
	}

	for _, table := range tables {
		attributeDefs := []types.AttributeDefinition{
			{AttributeName: aws.String(table.Key), AttributeType: types.ScalarAttributeTypeS},
		}
		keySchema := []types.KeySchemaElement{
			{AttributeName: aws.String(table.Key), KeyType: types.KeyTypeHash},
		}

		_, err := dynamoClient.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName:            aws.String(table.Name),
			AttributeDefinitions: attributeDefs,
			KeySchema:            keySchema,
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})

		if err != nil {
			fmt.Printf("Failed to create table %s: %v\n", table.Name, err)
		} else {
			fmt.Printf("Table %s created successfully.\n", table.Name)
		}

		if table.Name == "RebateClaim" {
			// Wait for the table to be active
			waitForTableToBeActive(ctx, dynamoClient, "RebateClaim")

			// Add the GSI (DateIndex)
			_, err = dynamoClient.UpdateTable(ctx, &dynamodb.UpdateTableInput{
				TableName: aws.String("RebateClaim"),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("Date"), AttributeType: types.ScalarAttributeTypeS},
				},
				GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{
					{
						Create: &types.CreateGlobalSecondaryIndexAction{
							IndexName: aws.String("DateIndex"),
							KeySchema: []types.KeySchemaElement{
								{AttributeName: aws.String("Date"), KeyType: types.KeyTypeHash}, // Date is the partition key for the GSI
							},
							Projection: &types.Projection{
								ProjectionType: types.ProjectionTypeAll, // Include all attributes
							},
							ProvisionedThroughput: &types.ProvisionedThroughput{
								ReadCapacityUnits:  aws.Int64(5),
								WriteCapacityUnits: aws.Int64(5),
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Printf("Failed to create GSI DateIndex for table %s: %v\n", table.Name, err)
			} else {
				fmt.Printf("GSI DateIndex created successfully for table %s.\n", table.Name)
			}
		}
	}

	// Additional uniqueness checks using GSIs
	addUniqueConstraint(ctx, dynamoClient, "RebateProgram", "ProgramName")
	addUniqueConstraint(ctx, dynamoClient, "RebateClaim", "TransactionID")
	addUniqueConstraint(ctx, dynamoClient, "Transaction", "ID")

	_, err := dynamoClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String("ReportCache"),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("CacheKey"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("CacheKey"), KeyType: types.KeyTypeHash},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		fmt.Printf("Failed to create ReportCache table: %v\n", err)
	} else {
		fmt.Println("ReportCache table created successfully.")
	}
}

func addUniqueConstraint(ctx context.Context, dynamoClient *dynamodb.Client, tableName, attributeName string) {
	_, err := dynamoClient.UpdateTable(ctx, &dynamodb.UpdateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String(attributeName), AttributeType: types.ScalarAttributeTypeS},
		},
		GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{
			{
				Create: &types.CreateGlobalSecondaryIndexAction{
					IndexName: aws.String(attributeName + "UniqueIndex"),
					KeySchema: []types.KeySchemaElement{
						{AttributeName: aws.String(attributeName), KeyType: types.KeyTypeHash},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeKeysOnly,
					},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(5),
						WriteCapacityUnits: aws.Int64(5),
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Failed to add unique constraint on %s in table %s: %v\n", attributeName, tableName, err)
	} else {
		fmt.Printf("Unique constraint added for %s in table %s.\n", attributeName, tableName)
	}
}

func waitForTableToBeActive(ctx context.Context, dynamoClient *dynamodb.Client, tableName string) {
	for {
		output, err := dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			fmt.Printf("Error describing table %s: %v\n", tableName, err)
			return
		}

		if output.Table.TableStatus == types.TableStatusActive {
			fmt.Printf("Table %s is now ACTIVE.\n", tableName)
			return
		}

		fmt.Printf("Table %s is still being created. Waiting...\n", tableName)
		time.Sleep(5 * time.Second)
	}
}
