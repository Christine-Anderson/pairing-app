package main

import (
	"lambda/api"
	"lambda/database"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleLambda(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db := database.NewDynamoDB()
	apiHandler := api.NewApiHandler(db)

	switch request.Path {
	case "/create-group":
		return apiHandler.CreateGroup(request)
	case "/join-group":
		return apiHandler.JoinGroup(request)
	case "/group-details/{groupId}":
		return apiHandler.GroupDetails(request)
	case "/match":
		return apiHandler.PerformMatching(request)
	default:
		return events.APIGatewayProxyResponse{
			Body:       "Not found",
			StatusCode: http.StatusNotFound,
		}, nil
	}
}

func main() {
	lambda.Start(HandleLambda)
}
