package api

import (
	"encoding/json"
	"lambda/database"
	"lambda/email"
	"lambda/types"

	"net/http"
	"net/mail"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	databaseStore database.DynamoDBClient
	emailService  email.EmailService
}

func NewApiHandler(databaseStore database.DynamoDBClient, emailService email.EmailService) ApiHandler {
	return ApiHandler{
		databaseStore: databaseStore,
		emailService:  emailService,
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

func verifyEmailSuccessResponse(email string) (events.APIGatewayProxyResponse, error) {
	responseBody := map[string]string{
		"email": email,
	}
	return successResponse(responseBody)
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

func (handler ApiHandler) VerifyEmail(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var groupDetails types.CreateGroupDetails

	err := json.Unmarshal([]byte(request.Body), &groupDetails)
	if err != nil {
		return errorResponse("Invalid Request: Error unmarshaling request body: "+err.Error(), http.StatusBadRequest), err
	}

	if !isValidEmail(groupDetails.Email) {
		return errorResponse("Invalid Request: Email is missing or invalid", http.StatusBadRequest), err
	}

	verifErr := handler.emailService.SendVerificationEmail(groupDetails.Email)
	if verifErr != nil {
		return errorResponse("Error verifying email: "+verifErr.Error(), http.StatusInternalServerError), err
	}

	return verifyEmailSuccessResponse(groupDetails.Email)
}

func (handler ApiHandler) CreateGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var groupDetails types.CreateGroupDetails

	err := json.Unmarshal([]byte(request.Body), &groupDetails)
	if err != nil {
		return errorResponse("Invalid Request: Error unmarshaling request body: "+err.Error(), http.StatusBadRequest), err
	}

	if groupDetails.Name == "" || groupDetails.GroupName == "" {
		return errorResponse("Invalid Request: Name or Group Name is missing", http.StatusBadRequest), err
	}

	groupId := uuid.New()
	newGroup := types.NewGroup(groupId.String(), groupDetails)
	dbErr := handler.databaseStore.AddGroup(newGroup)
	if dbErr != nil {
		return errorResponse("Error adding group to database: "+dbErr.Error(), http.StatusInternalServerError), err
	}

	emailErr := handler.emailService.SendConfirmationEmail(newGroup)
	if emailErr != nil {
		return errorResponse("Error sending confirmation email: "+emailErr.Error(), http.StatusInternalServerError), err
	}

	return groupSuccessResponse(groupId.String(), groupDetails.GroupName)
}

func (handler ApiHandler) JoinGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var groupDetails types.JoinGroupDetails
	err := json.Unmarshal([]byte(request.Body), &groupDetails)
	if err != nil {
		return errorResponse("Invalid Request: Error unmarshaling request body: "+err.Error(), http.StatusBadRequest), err
	}

	if groupDetails.Name == "" || !isValidUUID(groupDetails.GroupId) {
		return errorResponse("Invalid Request: Name or Group ID is missing or invalid", http.StatusBadRequest), err
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
