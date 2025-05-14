package provider

import "go.uber.org/fx"

var ProviderModule = fx.Module("provider_modules",
	fx.Provide(
		SetLogger,
		LoadViperConfig,
		GormPostgres,
		AwsConfig,
		NewRedisClient,
		NewCognitoClient,
		NewDynamoDBClient,
		NewTexTractClient,
		NewRekognitionClient,
		NewFiberApp,
	),
)
