package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// type MyEvent struct {
// 	Name string `json:"What is your name?"`
// 	Age  int    `json:"How old are you?"`
// }

// type MyResponse struct {
// 	Message string `json:"Answer"`
// }

// func HandleLambdaEvent(event *MyEvent) (*MyResponse, error) {
// 	if event == nil {
// 		return nil, fmt.Errorf("received nil event")
// 	}
// 	return &MyResponse{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil
// }

func createGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Group created successfully",
		StatusCode: http.StatusCreated,
	}, nil
}

func joinGroup(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Group joined successfully",
		StatusCode: http.StatusOK,
	}, nil
}

func groupDetails(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Group details",
		StatusCode: http.StatusOK,
	}, nil
}

func performMatching(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Matches",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/create-group":
			return createGroup(request)
		case "/join-group":
			return joinGroup(request)
		case "/group-details/{groupId}":
			return groupDetails(request)
		case "/match":
			return performMatching(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
