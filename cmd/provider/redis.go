package provider

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedisClient(v *viper.Viper, log *logrus.Logger) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     v.GetString("redis.addr"),
		Password: v.GetString("redis.password"),
		DB:       v.GetInt("redis.db"),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Error(err.Error())
	}
	return redisClient
}
