package main

import (
	"context"
	"fmt"
	"lambda/api"
	"lambda/database"
	"lambda/email"
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
	emailService, err := email.NewEmailService()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Email service error",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	apiHandler := api.NewApiHandler(db, emailService)

	switch request.Resource {
	case "/verify-email":
		return apiHandler.VerifyEmail(request)
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
