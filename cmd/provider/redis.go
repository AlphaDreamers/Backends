package provider

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"strconv"
)

func NewRedisClient(v *viper.Viper, logger *logrus.Logger) *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://default:jvRBJYdHdTXaKvhqVnuTPHLuCkwvypkL@maglev.proxy.rlwy.net:22038"
	}
	u, err := url.Parse(redisURL)
	if err != nil {
		logger.Fatalf("invalid redis URL: %v", err)
	}

	password, _ := u.User.Password()
	addr := u.Host
	db := 0

	if u.Path != "" && u.Path != "/" {
		dbStr := u.Path[1:]
		if parsedDB, err := strconv.Atoi(dbStr); err == nil {
			db = parsedDB
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Errorf("failed to connect to Redis: %v", err)
	}

	return redisClient
}
