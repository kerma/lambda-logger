package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/kerma/lambda-logger"
)

func HandleRequest(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log := logger.NewFromRequest(req)
	log.BindEnv()
	log.BindString("someKey", "test value")
	log.Println("All good")

	resp := events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           nil,
		MultiValueHeaders: nil,
		Body:              "",
		IsBase64Encoded:   false,
	}
	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
