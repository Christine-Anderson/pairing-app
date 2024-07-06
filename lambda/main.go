package main

import (
	"context"
	"fmt"
	"lambda/api"
	"lambda/database"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	db := database.NewDynamoDB("MyGroupTable")
	apiHandler := api.NewApiHandler(db)

	switch request.Resource {
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
	lambda.Start(handleRequest)
}
