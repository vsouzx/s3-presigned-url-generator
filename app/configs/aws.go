package configs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAWSConfig() (aws.Config) {
	config, err := config.LoadDefaultConfig(context.TODO());
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return config
}
