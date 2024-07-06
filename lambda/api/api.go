package api

import (
	"encoding/json"
	"lambda/database"
	"lambda/types"

	"net/http"
	"net/mail"

	"github.com/google/uuid"

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

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func errorResponse(message string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: statusCode,
	}
}

func createGroupResponse(groupId string, groupName string) (events.APIGatewayProxyResponse, error) {
	responseBody := map[string]string{
		"groupId":   groupId,
		"groupName": groupName,
	}

	jsonResponseBody, err := json.Marshal(responseBody)
	if err != nil {
		return errorResponse("Error creating group response", http.StatusInternalServerError), err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonResponseBody),
		StatusCode: http.StatusOK,
	}, nil
}

func groupDetailsResponse(group types.Group) (events.APIGatewayProxyResponse, error) {
	jsonResponseBody, err := json.Marshal(group)
	if err != nil {
		return errorResponse("Error creating group response", http.StatusInternalServerError), err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonResponseBody),
		StatusCode: http.StatusOK,
	}, nil
}

func (handler ApiHandler) CreateGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var groupDetails types.CreateGroupDetails

	err := json.Unmarshal([]byte(request.Body), &groupDetails)
	if err != nil {
		return errorResponse("Invalid Request: Error unmarshaling request body: "+err.Error(), http.StatusBadRequest), err
	}

	if groupDetails.Name == "" || !isValidEmail(groupDetails.Email) || groupDetails.GroupName == "" {
		return errorResponse("Invalid Request: Name, Email, or Group Name is missing or invalid", http.StatusBadRequest), err
	}

	groupId := uuid.New()
	newGroup := types.NewGroup(groupId.String(), groupDetails)
	dbErr := handler.databaseStore.AddGroup(newGroup)
	if dbErr != nil {
		return errorResponse("Error adding group to database: "+dbErr.Error(), http.StatusInternalServerError), err
	}

	// TODO send email with link to group page w/ JWT and groupId

	return createGroupResponse(groupId.String(), groupDetails.GroupName)
}

func (handler ApiHandler) JoinGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Group joined successfully",
		StatusCode: http.StatusOK,
	}, nil
}

func (handler ApiHandler) GroupDetails(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	groupId := request.PathParameters["groupId"]
	groupDetails, err := handler.databaseStore.FetchGroupById(groupId)
	if err != nil {
		return errorResponse("Error fetching group from database: "+err.Error(), http.StatusInternalServerError), err
	}

	return groupDetailsResponse(groupDetails)
}

func (handler ApiHandler) PerformMatching(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Matches",
		StatusCode: http.StatusOK,
	}, nil
}
