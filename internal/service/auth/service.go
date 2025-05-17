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
		clientId:           "7cgcs7ng3pbgd9elsn2f1lhund",
		clientSecret:       "1sa1nva5h3a0f11esvr4m051ui9u6khlnhjkipc5hgkqufiojdmf",
		textractClient:     textractClient,
		rekognitiionClient: rekognitiionClient,
		v:                  v,
		ctx:                context.Background(),
	}
}
