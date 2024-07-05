package api

import (
	"fmt"
	"lambda/database"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	databaseStore database.DynamoDBClient
}

func NewApiHandler(databaseStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		databaseStore: databaseStore,
	}
}

func (handler ApiHandler) CreateGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := handler.databaseStore.CreateGroupDB()
	fmt.Println(err)

	return events.APIGatewayProxyResponse{
		Body:       "Group created successfully",
		StatusCode: http.StatusCreated,
	}, nil
}

func (handler ApiHandler) JoinGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Group joined successfully",
		StatusCode: http.StatusOK,
	}, nil
}

func (handler ApiHandler) GroupDetails(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Group details",
		StatusCode: http.StatusOK,
	}, nil
}

func (handler ApiHandler) PerformMatching(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Matches",
		StatusCode: http.StatusOK,
	}, nil
}
