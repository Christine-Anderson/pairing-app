package util

import (
	"github.com/aws/aws-lambda-go/events"
)

func ErrorResponse(message string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: statusCode,
	}
}
