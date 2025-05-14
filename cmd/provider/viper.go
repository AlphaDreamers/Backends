package provider

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func LoadConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetDefault("app.env", "development")
	v.SetDefault("fiber.disableStartupMessage", false)
	v.SetDefault("fiber.prefork", false)
	v.SetDefault("fiber.caseSensitive", true)
	v.SetDefault("fiber.strictRouting", true)
	v.SetDefault("fiber.allowed-origin", "*")
	v.SetDefault("fiber.allowHeaders", "Origin, Content-Type, Accept, X-Csrf-Token")
	v.SetDefault("fiber.allowMethods", "GET, POST, PUT, DELETE, OPTIONS")
	v.SetDefault("fiber.allowCredentials", true)
	v.SetDefault("fiber.serverHeader", "auth-cognito")
	v.SetDefault("fiber.appName", "authWithCognito")
	v.SetDefault("fiber.port", "30011")
	v.SetDefault("fiber.idleTimeout", 120*time.Second)
	v.SetDefault("fiber.readTimeout", 10*time.Second)
	v.SetDefault("fiber.writeTimeout", 10*time.Second)
	v.SetDefault("limiter.max_requests", 100)
	v.SetDefault("limiter.expiration", 1*time.Minute)
	v.SetDefault("log.dir", "./logs")
	v.SetDefault("log.file", "./logs/app.log")
	v.SetDefault("log.max_size", 10)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 28)
	v.SetDefault("log.compress", true)
	v.SetDefault("log.pretty_print", true)
	v.SetDefault("client.timeout", 30*time.Second)
	v.SetDefault("aws.rds.local", "host=localhost user=postgres password=postgres dbname=auth sslmode=disable")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Config file not found, using defaults: %v\n", err)
	}

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	return v
}
