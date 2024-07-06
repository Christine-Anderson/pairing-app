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

func isValidUUID(groupId string) bool {
	if groupId == "" {
		return false
	}
	err := uuid.Validate(groupId)
	return err == nil
}

func errorResponse(message string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: statusCode,
	}
}

func groupSuccessResponse(groupId string, groupName string) (events.APIGatewayProxyResponse, error) {
	responseBody := map[string]string{
		"groupId":   groupId,
		"groupName": groupName,
	}
	return successResponse(responseBody)
}

func successResponse(responseBody any) (events.APIGatewayProxyResponse, error) {
	jsonResponseBody, err := json.Marshal(responseBody)
	if err != nil {
		return errorResponse("Error creating response", http.StatusInternalServerError), err
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

	return groupSuccessResponse(groupId.String(), groupDetails.GroupName)
}

func (handler ApiHandler) JoinGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var groupDetails types.JoinGroupDetails
	err := json.Unmarshal([]byte(request.Body), &groupDetails)
	if err != nil {
		return errorResponse("Invalid Request: Error unmarshaling request body: "+err.Error(), http.StatusBadRequest), err
	}
	// groupId, uuidErr := uuid.FromBytes([]byte(groupDetails.GroupId)) todo delete
	if groupDetails.Name == "" || !isValidEmail(groupDetails.Email) || !isValidUUID(groupDetails.GroupId) {
		return errorResponse("Invalid Request: Name, Email, or Group ID is missing or invalid", http.StatusBadRequest), err
	}

	groupToUpdate, err := handler.databaseStore.FetchGroupById(groupDetails.GroupId)
	if err != nil {
		return errorResponse("Error fetching group from database: "+err.Error(), http.StatusInternalServerError), err
	}

	groupToUpdate.GroupMembers = append(groupToUpdate.GroupMembers, types.NewGroupMember(groupDetails.Name, groupDetails.Email))

	dbErr := handler.databaseStore.UpdateGroup(groupToUpdate)
	if dbErr != nil {
		return errorResponse("Error updating group in database: "+dbErr.Error(), http.StatusInternalServerError), err
	}

	return groupSuccessResponse(groupToUpdate.GroupId, groupToUpdate.GroupName)
}

func (handler ApiHandler) GroupDetails(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	groupId := request.PathParameters["groupId"]
	groupDetails, err := handler.databaseStore.FetchGroupById(groupId)
	if err != nil {
		return errorResponse("Error fetching group from database: "+err.Error(), http.StatusInternalServerError), err
	}

	return successResponse(groupDetails)
}

func (handler ApiHandler) PerformMatching(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Matches",
		StatusCode: http.StatusOK,
	}, nil
}
