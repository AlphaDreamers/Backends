// main.go
package main

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/backend/test_mulit/provider"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
)

// LoadConfig sets up Viper
func LoadConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetDefault("app.env", "development")
	v.SetDefault("log.dir", "./logs")
	v.SetDefault("log.file", "./logs/app.log")
	v.SetDefault("log.max_size", 10)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 28)
	v.SetDefault("log.compress", true)
	v.SetDefault("log.pretty_print", true)
	v.SetDefault("server1.port", "30011")
	v.SetDefault("server1.name", "HealthServer")
	v.SetDefault("server1.header", "health-server")
	v.SetDefault("server1.read_timeout", "10s")
	v.SetDefault("server1.write_timeout", "10s")
	v.SetDefault("server1.idle_timeout", "120s")
	v.SetDefault("server2.port", "30012")
	v.SetDefault("server2.name", "PingServer")
	v.SetDefault("server2.header", "ping-server")
	v.SetDefault("server2.read_timeout", "10s")
	v.SetDefault("server2.write_timeout", "10s")
	v.SetDefault("server2.idle_timeout", "120s")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Config file not found, using defaults: %v\n", err)
	}

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	return v
}

// Server1Module defines the first Fiber server
var Server1Module = fx.Module("server1",
	fx.Provide(
		fx.Annotate(
			func(v *viper.Viper, logger *logrus.Logger) (*fiber.App, error) {
				app := provider.NewFiberApp(v, logger, "server1")
				app.Get("/health", func(c *fiber.Ctx) error {
					return c.JSON(fiber.Map{"status": "ok"})
				})
				return app, nil
			},
			fx.ResultTags(`name:"server1"`), // Tag to distinguish this *fiber.App
		),
	),
	fx.Invoke(
		fx.Annotate(
			func(app *fiber.App, v *viper.Viper, logger *logrus.Logger, lc fx.Lifecycle) {
				port := v.GetString("server1.port")
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						logger.Infof("Starting Server1 on :%s", port)
						go func() {
							if err := app.Listen(":" + port); err != nil {
								logger.Errorf("Server1 failed: %v", err)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("Shutting down Server1")
						return app.Shutdown()
					},
				})
			},
			fx.ParamTags(`name:"server1"`), // Match the tagged *fiber.App
		),
	),
)

// Server2Module defines the second Fiber server
var Server2Module = fx.Module("server2",
	fx.Provide(
		fx.Annotate(
			func(v *viper.Viper, logger *logrus.Logger) (*fiber.App, error) {
				app := provider.NewFiberApp(v, logger, "server2")
				app.Get("/ping", func(c *fiber.Ctx) error {
					return c.JSON(fiber.Map{"message": "pong"})
				})
				return app, nil
			},
			fx.ResultTags(`name:"server2"`), // Tag to distinguish this *fiber.App
		),
	),
	fx.Invoke(
		fx.Annotate(
			func(app *fiber.App, v *viper.Viper, logger *logrus.Logger, lc fx.Lifecycle) {
				port := v.GetString("server2.port")
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						logger.Infof("Starting Server2 on :%s", port)
						go func() {
							if err := app.Listen(":" + port); err != nil {
								logger.Errorf("Server2 failed: %v", err)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("Shutting down Server2")
						return app.Shutdown()
					},
				})
			},
			fx.ParamTags(`name:"server2"`), // Match the tagged *fiber.App
		),
	),
)

func main() {
	app := fx.New(
		fx.Provide(LoadConfig),
		fx.Provide(provider.SetLogger),
		Server1Module,
		Server2Module,
		fx.Invoke(func(logger *logrus.Logger, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info("Application started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Application stopped")
					return nil
				},
			})
		}),
	)

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		if err := app.Stop(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to stop app: %v\n", err)
		}
	}()

	// Start the application
	if err := app.Start(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start app: %v\n", err)
		os.Exit(1)
	}

	// Wait for the application to stop
	<-app.Done()
}
