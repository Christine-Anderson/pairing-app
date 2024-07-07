package database

import (
	"errors"
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
		TableName: aws.String(client.tableName),
		Item:      item,
	}

	_, err = client.databaseStore.PutItem(input)
	if err != nil {
		return fmt.Errorf("error putting item: %w", err)

	}

	return nil
}

func (client DynamoDBClient) FetchGroupById(groupId string) (types.Group, error) {
	var group types.Group

	result, err := client.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(client.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"groupId": {
				S: aws.String(groupId),
			},
		},
	})

	if err != nil {
		return group, fmt.Errorf("error getting item: %w", err)
	}

	if result.Item == nil {
		return group, errors.New("could not find group ID " + groupId)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &group)
	if err != nil {
		return group, fmt.Errorf("error marshalling item: %w", err)
	}
	return group, nil
}

func (client DynamoDBClient) AddGroupMember(groupId string, newMemberName string, newMemberEmail string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(client.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"groupId": {
				S: aws.String(groupId),
			},
		},
		UpdateExpression: aws.String("set groupMembers = list_append(groupMembers, :m)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":m": {
				L: []*dynamodb.AttributeValue{
					{M: map[string]*dynamodb.AttributeValue{
						"name":  {S: aws.String(newMemberName)},
						"email": {S: aws.String(newMemberEmail)},
					}},
				},
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err := client.databaseStore.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("error updating item: %w", err)
	}

	return nil
}
