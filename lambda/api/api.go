package api

import (
	"encoding/json"
	"fmt"
	"lambda/assignmentGenerator"
	"lambda/database"
	"lambda/email"
	"lambda/jwt"
	"lambda/types"
	"lambda/util"
	"log"

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
		return util.ErrorResponse("Error creating response"+err.Error(), http.StatusInternalServerError), err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonResponseBody),
		StatusCode: http.StatusOK,
	}, nil
}

func (handler ApiHandler) VerifyEmail(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var groupDetails types.CreateGroupDetails

	jsonErr := json.Unmarshal([]byte(request.Body), &groupDetails)
	if jsonErr != nil {
		return util.ErrorResponse("Invalid Request: Error unmarshaling request body: "+jsonErr.Error(), http.StatusBadRequest), jsonErr
	}

	if !isValidEmail(groupDetails.Email) {
		return util.ErrorResponse("Invalid Request: Email is missing or invalid", http.StatusBadRequest), fmt.Errorf("email is missing or invalid")
	}

	verifErr := handler.emailService.SendVerificationEmail(groupDetails.Email)
	if verifErr != nil {
		return util.ErrorResponse("Error verifying email: "+verifErr.Error(), http.StatusInternalServerError), verifErr
	}

	return verifyEmailSuccessResponse(groupDetails.Email)
}

func (handler ApiHandler) CreateGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var groupDetails types.CreateGroupDetails

	jsonErr := json.Unmarshal([]byte(request.Body), &groupDetails)
	if jsonErr != nil {
		return util.ErrorResponse("Invalid Request: Error unmarshaling request body: "+jsonErr.Error(), http.StatusBadRequest), jsonErr
	}

	if groupDetails.Name == "" || groupDetails.GroupName == "" {
		return util.ErrorResponse("Invalid Request: Name or Group Name is missing", http.StatusBadRequest), fmt.Errorf("name or group name is missing")
	}

	groupId := uuid.New()
	memberId := uuid.New()
	newGroup := types.NewGroup(groupId.String(), memberId.String(), groupDetails)
	dbErr := handler.databaseStore.AddGroup(newGroup)
	if dbErr != nil {
		return util.ErrorResponse("Error adding group to database: "+dbErr.Error(), http.StatusInternalServerError), dbErr
	}

	emailErr := handler.emailService.SendConfirmationEmail(newGroup)
	if emailErr != nil {
		return util.ErrorResponse("Error sending confirmation email: "+emailErr.Error(), http.StatusInternalServerError), emailErr
	}

	return groupSuccessResponse(groupId.String(), groupDetails.GroupName)
}

func (handler ApiHandler) JoinGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var groupDetails types.JoinGroupDetails
	jsonErr := json.Unmarshal([]byte(request.Body), &groupDetails)
	if jsonErr != nil {
		return util.ErrorResponse("Invalid Request: Error unmarshaling request body: "+jsonErr.Error(), http.StatusBadRequest), jsonErr
	}

	if groupDetails.Name == "" || !isValidUUID(groupDetails.GroupId) {
		return util.ErrorResponse("Invalid Request: Name or Group ID is missing or invalid", http.StatusBadRequest), fmt.Errorf("name or group ID is missing or invalid")
	}

	joinedGroup, dbFetchErr := handler.databaseStore.FetchGroupById(groupDetails.GroupId)
	if dbFetchErr != nil {
		return util.ErrorResponse("Error fetching group from database: "+dbFetchErr.Error(), http.StatusInternalServerError), dbFetchErr
	}

	memberId := uuid.New()
	dbUpdateErr := handler.databaseStore.AddGroupMember(groupDetails.GroupId, memberId.String(), groupDetails.Name, groupDetails.Email)
	if dbUpdateErr != nil {
		return util.ErrorResponse("Error updating group in database: "+dbUpdateErr.Error(), http.StatusInternalServerError), dbUpdateErr
	}

	return groupSuccessResponse(groupDetails.GroupId, joinedGroup.GroupName)
}

func (handler ApiHandler) GroupDetails(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	groupId := request.PathParameters["groupId"]
	tokenString := request.QueryStringParameters["jwt"]

	tokenErr := jwt.VerifyToken(tokenString, groupId)
	if tokenErr != nil {
		return util.ErrorResponse("Unauthorized: "+tokenErr.Error(), http.StatusUnauthorized), tokenErr
	}

	groupDetails, dbFetchErr := handler.databaseStore.FetchGroupById(groupId)
	if dbFetchErr != nil {
		return util.ErrorResponse("Error fetching group from database: "+dbFetchErr.Error(), http.StatusInternalServerError), dbFetchErr
	}

	return successResponse(groupDetails)
}

func (handler ApiHandler) GenerateAssignments(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	groupId := request.PathParameters["groupId"]
	tokenString := request.QueryStringParameters["jwt"]

	tokenErr := jwt.VerifyToken(tokenString, groupId)
	if tokenErr != nil {
		return util.ErrorResponse("Unauthorized: "+tokenErr.Error(), http.StatusUnauthorized), tokenErr
	}

	var restrictions types.AssignmentRestrictions
	jsonErr := json.Unmarshal([]byte(request.Body), &restrictions)
	if jsonErr != nil {
		return util.ErrorResponse("Invalid Request: Error unmarshaling request body: "+jsonErr.Error(), http.StatusBadRequest), jsonErr
	}

	log.Printf("Received assignment restrictions: %s", restrictions)

	groupDetails, dbFetchErr := handler.databaseStore.FetchGroupById(groupId)
	if dbFetchErr != nil {
		return util.ErrorResponse("Error fetching group from database: "+dbFetchErr.Error(), http.StatusInternalServerError), dbFetchErr
	}

	assignments, assignErr := assignmentGenerator.GenerateAssignments(groupDetails, restrictions.Restrictions)
	if assignErr != nil {
		return util.ErrorResponse("Error generating assignments: "+assignErr.Error(), http.StatusBadRequest), assignErr
	}

	handler.emailService.SendAssignmentEmails(assignments, groupDetails)

	return events.APIGatewayProxyResponse{
		Body:       "Assignments generated successfully. Group members have been emailed the results.",
		StatusCode: http.StatusOK,
	}, nil
}
