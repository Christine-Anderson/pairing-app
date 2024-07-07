package main

import (
	"context"
	"fmt"
	"lambda/api"
	"lambda/database"
	"lambda/email"
	"lambda/util"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	log.Printf("Body = %s.\n", request.Body)
	log.Println("Headers:")
	for key, value := range request.Headers {
		log.Printf("    %s: %s\n", key, value)
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
	case "/assign/{groupId}":
		return apiHandler.GenerateAssignments(request)
	default:
		return util.ErrorResponse("Resource Not found", http.StatusNotFound), fmt.Errorf("resource not found")
	}
}

func main() {
	lambda.Start(handleRequest)
}
