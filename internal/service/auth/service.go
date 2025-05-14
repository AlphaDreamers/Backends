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
		clientId:           "7qllcjjcq7p506kq88vkfiu92g",
		clientSecret:       "1ipuga7399127snjbbgletfpr25lk6hleucb5fptn6nvrefn40ri",
		textractClient:     textractClient,
		rekognitiionClient: rekognitiionClient,
		v:                  v,
		ctx:                context.Background(),
	}
}
