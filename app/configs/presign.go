package configs

import "github.com/aws/aws-sdk-go-v2/service/s3"

type PresignConfig struct {
	PresignClient *s3.PresignClient
}

func NewPresignConfig(client *s3.PresignClient) *PresignConfig {
	return &PresignConfig{
		PresignClient: client,
	}
}