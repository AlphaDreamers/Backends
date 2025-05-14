package provider

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func AwsConfig(v *viper.Viper, log *logrus.Logger) *aws.Config {
	region := v.GetString("aws.region")
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region))
	if err != nil {
		log.WithError(err).Fatal("failed to load aws config")
	}
	return &cfg
}

func NewCognitoClient(
	awsConfig *aws.Config,
) *cognitoidentityprovider.Client {
	client := cognitoidentityprovider.NewFromConfig(*awsConfig)
	return client
}

func NewS3Client(
	awsConfig *aws.Config,
) *s3.Client {
	client := s3.NewFromConfig(*awsConfig)
	return client
}

func NewDynamoDBClient(
	awsConfig *aws.Config) *dynamodb.Client {
	client := dynamodb.NewFromConfig(*awsConfig)
	return client
}

func NewTexTractClient(
	awsConfig *aws.Config) *textract.Client {
	client := textract.NewFromConfig(*awsConfig)
	return client
}

func NewRekognitionClient(
	awsConfig *aws.Config) *rekognition.Client {
	client := rekognition.NewFromConfig(*awsConfig)
	return client
}
