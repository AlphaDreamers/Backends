package provider

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func AwsConfig(v *viper.Viper, log *logrus.Logger) *aws.Config {
	region := v.GetString("aws.region")
	accessKey := os.Getenv("ACESS_KEY_ID")
	secretKey := os.Getenv("SECRECT_ACCESS_KEY")
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			)))
	if err != nil {
		log.Fatal(err.Error())
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
