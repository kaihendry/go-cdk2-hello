package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// respond with hello world

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "Hello World",
	}, nil
}
