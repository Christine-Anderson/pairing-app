package api

import (
	"encoding/json"
	"fmt"
	"lambda/database"
	"lambda/email"
	"lambda/types"
	"lambda/util"
	"log"
	"math/rand"
	"strings"

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
	newGroup := types.NewGroup(groupId.String(), groupDetails)
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

	dbUpdateErr := handler.databaseStore.AddGroupMember(groupDetails.GroupId, groupDetails.Name, groupDetails.Email)
	if dbUpdateErr != nil {
		return util.ErrorResponse("Error updating group in database: "+dbUpdateErr.Error(), http.StatusInternalServerError), dbUpdateErr
	}

	return groupSuccessResponse(groupDetails.GroupId, joinedGroup.GroupName)
}

func (handler ApiHandler) GroupDetails(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	groupId := request.PathParameters["groupId"]
	tokenString := request.QueryStringParameters["jwt"]

	tokenErr := util.VerifyToken(tokenString, groupId)
	if tokenErr != nil {
		return util.ErrorResponse("Unauthorized: "+tokenErr.Error(), http.StatusUnauthorized), tokenErr
	}

	groupDetails, dbFetchErr := handler.databaseStore.FetchGroupById(groupId)
	if dbFetchErr != nil {
		return util.ErrorResponse("Error fetching group from database: "+dbFetchErr.Error(), http.StatusInternalServerError), dbFetchErr
	}

	return successResponse(groupDetails)
}

func getDisallowedAssignments(assignments map[string]string, restrictions string) map[string][]string {
	disallowedAssignments := make(map[string][]string)

	parts := strings.Split(restrictions, ";")

	for _, part := range parts {
		disallowedAssignment := strings.Split(part, ",")
		if len(disallowedAssignment) == 2 {
			giver := strings.TrimSpace(disallowedAssignment[0])
			receiver := strings.TrimSpace(disallowedAssignment[1])

			_, ok0 := assignments[giver]
			_, ok1 := assignments[receiver]
			if ok0 && ok1 {
				log.Printf("Adding restriction: %s cannot give to %s\n", giver, receiver)
				disallowedAssignments[giver] = append(disallowedAssignments[giver], receiver)
			}
		}
	}

	log.Printf("Disallowed assignments calculated: %s", disallowedAssignments)
	return disallowedAssignments
}

func isAssignmentAllowed(giver string, receiver string, disallowedAssignments map[string][]string) bool {
	value, ok := disallowedAssignments[giver]
	if (ok && util.Contains(value, receiver)) || giver == receiver {
		log.Printf("Assignment from %s to %s is not allowed\n", giver, receiver)
		return false
	}
	log.Printf("Assignment from %s to %s is allowed\n", giver, receiver)
	return true
}

func isAssignmentComplete(assignments map[string]string) bool {
	for _, value := range assignments {
		if value == "" {
			log.Printf("Assignments are not complete: %s", assignments)
			return false
		}
	}
	log.Printf("Assignments are complete: %s", assignments)
	return true
}

func generateCandidateAssignments(currAssignment map[string]string, disallowedAssignments map[string][]string, giver string, unassignedReceivers []string) []map[string]string {
	var candidateAssignments []map[string]string

	for _, receiver := range unassignedReceivers {
		if isAssignmentAllowed(giver, receiver, disallowedAssignments) {
			candidate := util.Copy(currAssignment)
			candidate[giver] = receiver
			candidateAssignments = append(candidateAssignments, candidate)
		}
	}

	rand.Shuffle(len(candidateAssignments), func(i, j int) {
		candidateAssignments[i], candidateAssignments[j] = candidateAssignments[j], candidateAssignments[i]
	})

	log.Printf("Candidate assignments generated: %s", candidateAssignments)
	return candidateAssignments
}

func backtrackAssign(currAssignment map[string]string, disallowedAssignments map[string][]string, unassignedReceivers []string) bool {
	if isAssignmentComplete(currAssignment) {
		return true
	}

	giver := ""
	for key, value := range currAssignment {
		if value == "" {
			giver = key
			break
		}
	}

	candidateAssignments := generateCandidateAssignments(currAssignment, disallowedAssignments, giver, unassignedReceivers)
	for _, candidateAssignments := range candidateAssignments {
		unassignedReceivers = util.Remove(unassignedReceivers, candidateAssignments[giver])
		currAssignment[giver] = candidateAssignments[giver]
		if backtrackAssign(currAssignment, disallowedAssignments, unassignedReceivers) {
			return true
		}
		unassignedReceivers = append(unassignedReceivers, candidateAssignments[giver])
		currAssignment[giver] = ""
	}

	return false
}

func getAssignments(group types.Group, restrictions string) (map[string]string, error) {
	assignments := make(map[string]string)

	var unassignedReceivers []string
	for _, member := range group.GroupMembers {
		assignments[member.Name] = ""
		unassignedReceivers = append(unassignedReceivers, member.Name)
	}

	disallowedAssignments := getDisallowedAssignments(assignments, restrictions)

	if backtrackAssign(assignments, disallowedAssignments, unassignedReceivers) {
		return assignments, nil
	} else {
		return assignments, fmt.Errorf("no valid assignments")
	}
}

func (handler ApiHandler) GenerateAssignments(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	groupId := request.PathParameters["groupId"]
	tokenString := request.QueryStringParameters["jwt"]

	tokenErr := util.VerifyToken(tokenString, groupId)
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

	assignments, assignErr := getAssignments(groupDetails, restrictions.Restrictions)
	if assignErr != nil {
		return util.ErrorResponse("Error generating assignments: "+assignErr.Error(), http.StatusBadRequest), assignErr
	}

	// handler.emailService.SendAssignmentEmails(assignments) todo
	log.Printf("Assignments generated successfully: %s", assignments)

	return events.APIGatewayProxyResponse{
		Body:       "test",
		StatusCode: http.StatusOK,
	}, nil
}
