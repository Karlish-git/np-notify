package models

import (
	"context"
	"np-notify/pkg/aws_helpers"
	"np-notify/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Account struct {
	Uuid         string                 `dynamodbav:"uuid"`
	Username     string                 `dynamodbav:"username"`
	userId       string                 `dynamodbav:"userId"`
	Email        string                 `dynamodbav:"email"`
	Token        string                 `dynamodbav:"token"`
	UnreadCounts map[string]UnreadCount `dynamodbav:"unreads"`
}

func (a *Account) getKey() map[string]types.AttributeValue {
	uuid, err := attributevalue.Marshal(a.Uuid)
	if err != nil {
		utils.LogError("Failed to marshal uuid: %v", err)
	}
	return map[string]types.AttributeValue{"uuid": uuid}

}

func GetAllAccounts(ctx context.Context, db aws_helpers.TableBasics) ([]Account, error) {
	var acs []Account
	var err error
	paginator := dynamodb.NewScanPaginator(db.DynamoDbClient, &dynamodb.ScanInput{
		TableName: aws.String(db.TableName),
	})

	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			utils.LogError("Account paginator Filed with error %v", err)
		} else {
			var accPage []Account
			err = attributevalue.UnmarshalListOfMaps(response.Items, &accPage)
			if err != nil {
				utils.LogError("Failed to unmarshal account page: %v", err)
			} else {
				acs = append(acs, accPage...)
			}
		}
	}
	// fix issue of nil Inreadcounts
	return acs, err
}

func (a *Account) UpdateUnreads(ctx context.Context, db aws_helpers.TableBasics) error {
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]UnreadCount

	update := expression.Set(expression.Name("unreads"), expression.Value(a.UnreadCounts))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		utils.LogError("Failed to build update expression: %v", err)
		return err
	} else {
		response, err = db.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(db.TableName),
			Key:                       a.getKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			utils.LogError("Couldn't update unreads for account %v. Here's why: %v\n", a.Username, err)
			return err
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				utils.LogError("Couldn't unmarshal update response. Here's why: %v\n", err)
				return err
			}
		}
	}
	return err
}
