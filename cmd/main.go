// main.go
package main

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/backend/cmd/auth"
	"github.com/SwanHtetAungPhyo/backend/cmd/chat"
	"github.com/SwanHtetAungPhyo/backend/cmd/gig"
	"github.com/SwanHtetAungPhyo/backend/cmd/order"
	"github.com/SwanHtetAungPhyo/backend/cmd/provider"
	"github.com/SwanHtetAungPhyo/backend/cmd/user"
	"github.com/SwanHtetAungPhyo/backend/cmd/wallet"
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
	// Defaults for each module
	modules := []string{"auth", "chat", "gig", "wallet", "order", "user"}
	for i, mod := range modules {
		v.SetDefault(mod+".port", fmt.Sprintf("300%02d", 11+i))
		v.SetDefault(mod+".name", mod+"Server")
		v.SetDefault(mod+".header", mod+"-server")
		v.SetDefault(mod+".read_timeout", "10s")
		v.SetDefault(mod+".write_timeout", "10s")
		v.SetDefault(mod+".idle_timeout", "120s")
		v.SetDefault(mod+".certificate.cert", "/certificates/cert.pem")
		v.SetDefault(mod+".certificate.key", "/certificates/key.pem")
	}

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Config file not found, using defaults: %v\n", err)
	}

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	return v
}

func main() {
	app := fx.New(
		fx.Provide(LoadConfig),
		fx.Provide(provider.SetLogger),
		auth.ServerStateModule,
		chat.ServerStateModule,
		gig.ServerStateModule,
		wallet.ServerStateModule,
		order.ServerStateModule,
		user.ServerStateModule,
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
			_, err2 := fmt.Fprintf(os.Stderr, "Failed to stop app: %v\n", err)
			if err2 != nil {
				fmt.Printf("Failed to stop app: %v\n", err.Error())
				return
			}
		}
	}()

	if err := app.Start(context.Background()); err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, "Failed to start app: %v\n", err)
		if err2 != nil {
			fmt.Printf("Failed to start app: %v\n", err.Error())
			return
		}
		os.Exit(1)
	}

	<-app.Done()
}
