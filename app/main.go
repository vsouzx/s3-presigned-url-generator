package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vsouzx/s3-presigned-url-generator/configs"
	"github.com/vsouzx/s3-presigned-url-generator/service"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	presignConfig := configs.NewPresignConfig(s3.NewPresignClient(s3.NewFromConfig(configs.GetAWSConfig())))
	presignClient := presignConfig.PresignClient

	service := service.NewGeneratePresignedURLService(presignClient)
	return service.GeneratePresignedURL(ctx, req)
}
