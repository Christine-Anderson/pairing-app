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

func addCorsHeaders(response events.APIGatewayProxyResponse) events.APIGatewayProxyResponse {
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	response.Headers["Access-Control-Allow-Origin"] = "*"
	response.Headers["Access-Control-Allow-Headers"] = "Content-Type"
	response.Headers["Access-Control-Allow-Methods"] = "OPTIONS,GET,POST"
	return response
}

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
		return addCorsHeaders(util.ErrorResponse("Internal Server Error: "+esErr.Error(), http.StatusInternalServerError)), esErr
	}
	apiHandler := api.NewApiHandler(db, emailService)

	var response events.APIGatewayProxyResponse
	var err error

	switch request.Resource {
	case "/verify-email":
		response, err = apiHandler.VerifyEmail(request)
	case "/create-group":
		response, err = apiHandler.CreateGroup(request)
	case "/join-group":
		response, err = apiHandler.JoinGroup(request)
	case "/group-details/{groupId}":
		response, err = apiHandler.GroupDetails(request)
	case "/assign/{groupId}":
		response, err = apiHandler.GenerateAssignments(request)
	default:
		return addCorsHeaders(util.ErrorResponse("Resource Not found", http.StatusNotFound)), fmt.Errorf("resource not found")
	}

	response = addCorsHeaders(response)
	return response, err
}

func main() {
	lambda.Start(handleRequest)
}
