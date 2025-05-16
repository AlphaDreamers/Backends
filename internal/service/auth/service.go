package auth

import (
	"context"
	"github.com/SwanHtetAungPhyo/backend/internal/repo/auth"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Service struct {
	repo               *auth.Repository
	log                *logrus.Logger
	cognitoClient      *cognitoidentityprovider.Client
	clientId           string
	clientSecret       string
	textractClient     *textract.Client
	rekognitiionClient *rekognition.Client
	v                  *viper.Viper
	ctx                context.Context
}

func NewService(
	repo *auth.Repository,
	log *logrus.Logger,
	cognitoClient *cognitoidentityprovider.Client,
	textractClient *textract.Client,
	rekognitiionClient *rekognition.Client,
	v *viper.Viper,
) *Service {
	return &Service{
		repo:               repo,
		log:                log,
		cognitoClient:      cognitoClient,
		clientId:           "799269a72q2ih238tm10dveqpk",
		clientSecret:       "o5v3vvjs4fkstemagnagjiu0bii4le943v0gfviohig3utdqkh3",
		textractClient:     textractClient,
		rekognitiionClient: rekognitiionClient,
		v:                  v,
		ctx:                context.Background(),
	}
}
