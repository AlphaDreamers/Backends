package main

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/backend/cmd/auth"
	"github.com/SwanHtetAungPhyo/backend/cmd/provider"
	ah "github.com/SwanHtetAungPhyo/backend/internal/handler/auth"
	ar "github.com/SwanHtetAungPhyo/backend/internal/repo/auth"
	as "github.com/SwanHtetAungPhyo/backend/internal/service/auth"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
)

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
	modules := []string{"auth", "chat", "gig", "wallet", "order", "user"}
	for i, mod := range modules {
		v.SetDefault(mod+".port", fmt.Sprintf("300%02d", 11+i))
		v.SetDefault(mod+".name", mod+"Server")
		v.SetDefault(mod+".header", mod+"-server")
		v.SetDefault(mod+".read_timeout", "10s")
		v.SetDefault(mod+".write_timeout", "10s")
		v.SetDefault(mod+".idle_timeout", "120s")
		v.SetDefault(mod+".certificates.cert", "/certificates/cert.pem")
		v.SetDefault(mod+".certificates.key", "/certificates/key.pem")
	}

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Config file not found, using defaults: %v\n", err)
	}

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	return v
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err.Error())
	}
	app := fx.New(
		fx.Provide(LoadConfig),
		provider.ProviderModule,
		fx.Provide(as.NewService),
		fx.Provide(ar.NewRepository),
		fx.Provide(ah.NewHandler),
		auth.ServerStateModule,
		//chat.ServerStateModule,
		//gig.ServerStateModule,
		//wallet.ServerStateModule,
		//order.ServerStateModule,
		//user.ServerStateModule,
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		if err := app.Stop(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to stop app: %v\n", err)
		}
	}()

	if err := app.Start(context.Background()); err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, "Failed to start app: %v\n", err)
		if err2 != nil {
			return
		}
		os.Exit(1)
	}

	<-app.Done()
}
