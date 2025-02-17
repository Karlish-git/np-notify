package aws_helpers

import (
	"context"
	"np-notify/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func LocalAWS(
	ctx context.Context, tableName string) (TableBasics, SnsActions, error) {
	awsEndpoint := "http://localhost:4566"
	awsRegion := "eu-north-1"

	var err error
	var db TableBasics
	var snsActions SnsActions

	// awsCfg, err := config.LoadDefaultConfig(ctx)
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return db, snsActions, err
	}

	dynamoClient := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	db = TableBasics{
		DynamoDbClient: dynamoClient,
		TableName:      tableName,
	}
	_, err = db.TableExists(ctx)
	if err != nil {
		utils.LogError("TABLE DOES NOT EXIST")
		panic(err)
	}

	snsClient := sns.NewFromConfig(awsCfg, func(o *sns.Options) {
		o.BaseEndpoint = aws.String(awsEndpoint)
	})
	snsActions = SnsActions{SnsClient: snsClient}

	return db, snsActions, nil
}

func AWS(
	ctx context.Context, tableName string) (TableBasics, SnsActions, error) {
	var err error
	var db TableBasics
	var snsActions SnsActions

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return db, snsActions, err
	}

	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	db = TableBasics{
		DynamoDbClient: dynamoClient,
		TableName:      tableName,
	}
	_, err = db.TableExists(ctx)
	if err != nil {
		utils.LogError("TABLE DOES NOT EXIST")
		panic(err)
	}

	snsClient := sns.NewFromConfig(awsCfg)
	snsActions = SnsActions{SnsClient: snsClient}

	return db, snsActions, nil
}
