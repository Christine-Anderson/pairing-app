package database

import (
	"fmt"
	"lambda/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
	tableName     string
}

func NewDynamoDB(tableName string) DynamoDBClient {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	return DynamoDBClient{
		databaseStore: svc,
		tableName:     tableName,
	}
}

func (client DynamoDBClient) AddGroup(newGroup types.Group) error {
	item, err := dynamodbattribute.MarshalMap(newGroup)
	if err != nil {
		return fmt.Errorf("error marshalling group: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(client.tableName),
	}

	_, err = client.databaseStore.PutItem(input)
	if err != nil {
		return fmt.Errorf("error putting item into DynamoDB: %w", err)

	}

	return nil
}
