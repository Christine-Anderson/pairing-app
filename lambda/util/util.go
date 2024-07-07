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

func Contains(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}

func Remove(slice []string, s string) []string {
	for i, value := range slice {
		if value == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func Copy(m map[string]string) map[string]string {
	newMap := make(map[string]string)
	for key, value := range m {
		m[key] = value
	}
	return newMap
}
