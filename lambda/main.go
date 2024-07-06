package main

import (
	"context"
	"fmt"
	"lambda/api"
	"lambda/database"
	"lambda/email"
	"lambda/util"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	envErr := godotenv.Load()
	if envErr != nil {
		return util.ErrorResponse("Internal Server Error: "+envErr.Error(), http.StatusInternalServerError), envErr
	}

	db := database.NewDynamoDB("MyGroupTable")
	emailService, esErr := email.NewEmailService()

	if esErr != nil {
		return util.ErrorResponse("Internal Server Error: "+esErr.Error(), http.StatusInternalServerError), esErr
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
		return util.ErrorResponse("Resource Not found", http.StatusNotFound), fmt.Errorf("Resource Not Found")
	}
}

func main() {
	lambda.Start(handleRequest)
}
