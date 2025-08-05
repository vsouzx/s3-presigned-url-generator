package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vsouzx/s3-presigned-url-generator/dto"
)

type GeneratePresignedURLService struct {
	presignClient *s3.PresignClient
}

func NewGeneratePresignedURLService(presignClient *s3.PresignClient) *GeneratePresignedURLService {
	return &GeneratePresignedURLService{
		presignClient: presignClient,
	}
}	

func (gpu *GeneratePresignedURLService)GeneratePresignedURL(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucketName := os.Getenv("BUCKET_NAME")
	fmt.Println("Bucket: " + bucketName)

	fileName := req.QueryStringParameters["fileName"]
	if fileName == "" {
		fmt.Println("fileName parameter is missing")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "Missing fileName param"}, nil
	}

	reqPut := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}

	presignedReq, err := gpu.presignClient.PresignPutObject(ctx, reqPut, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		fmt.Println("Error generating presigned URL:", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error ao gerar URL pre assinada: %v", err)}, nil
	}

	body, _ := json.Marshal(dto.Response{URL: presignedReq.URL})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
