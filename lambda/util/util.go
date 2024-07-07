package util

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

type customClaims struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	GroupId string `json:"groupId"`
	jwt.RegisteredClaims
}

func ErrorResponse(message string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: statusCode,
	}
}

const (
	NUM_DAYS_UNTIL_EXPIRY = 24 * 30 * 3
)

func getSecretKey() []byte {
	secret := os.Getenv("SECRET_KEY")
	return []byte(secret)
}

func CreateToken(groupId string) (string, error) {
	secretKey := getSecretKey()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"groupId": groupId,
			"exp":     time.Now().Add(time.Hour * NUM_DAYS_UNTIL_EXPIRY).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, groupId string) error {
	if tokenString == "" {
		return fmt.Errorf("authorization token required")
	}

	secretKey := getSecretKey()

	token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if parseErr != nil {
		return parseErr
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return fmt.Errorf("invalid token")
	}
	expiration := time.Unix(int64(claims["exp"].(float64)), 0)
	if expiration.Before(time.Now()) {
		return fmt.Errorf("expired token")
	}

	tokenGroupId := claims["groupId"].(string)
	if tokenGroupId != groupId {
		return fmt.Errorf("invalid token")
	}

	return nil
}
