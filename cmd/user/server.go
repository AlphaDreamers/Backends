package user

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// ServerStateModule defines the Fx module for the user server
var ServerStateModule = fx.Module("user",
	fx.Provide(
		fx.Annotate(
			NewServerState,
			fx.ResultTags(`name:"user"`), // Tag to distinguish this *fiber.App
		),
	),
	fx.Invoke(
		fx.Annotate(
			RegisterLifeCycle,
			fx.ParamTags(`name:"user"`), // Match the tagged *fiber.App
		),
	),
)

type ServerState struct {
	log      *logrus.Logger
	fiberApp *fiber.App
	v        *viper.Viper
}

func NewServerState(
	log *logrus.Logger,
	v *viper.Viper,
) *ServerState {
	fiberApp := providerr.NewFiberApp(v, log, "user")
	// Define a basic endpoint for the user server
	fiberApp.Get("/user/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "user server running"})
	})
	return &ServerState{
		log:      log,
		fiberApp: fiberApp,
		v:        v,
	}
}

func (s *ServerState) Start() error {
	pwd, _ := os.Getwd()
	cert := pwd + s.v.GetString("user.certificate.cert")
	key := pwd + s.v.GetString("user.certificate.key")
	port := s.v.GetString("user.port")

	s.log.Infof("Starting user server on port %s...", port)

	go func() {
		if s.v.GetString("app.env") == "production" {
			if err := s.fiberApp.ListenTLS(":"+port, cert, key); err != nil {
				s.log.Errorf("Failed to start user server: %v", err)
			}
		} else {
			if err := s.fiberApp.Listen(":" + port); err != nil {
				s.log.Errorf("Failed to start user server: %v", err)
			}
		}
	}()
	return nil
}

func (s *ServerState) Stop() error {
	s.log.Infof("Shutting down user server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		s.log.Errorf("Failed to gracefully shut down user server: %v", err)
		return err
	}

	s.log.Infof("user server shutdown complete.")
	return nil
}

// RegisterLifeCycle registers the lifecycle with Fx
func RegisterLifeCycle(lc fx.Lifecycle, s *ServerState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Infof("Starting user lifecycle...")
			if err := s.Start(); err != nil {
				s.log.Errorf("Failed to start user lifecycle: %v", err)
				return err
			}
			s.log.Infof("user lifecycle started.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Infof("Stopping user lifecycle...")
			if err := s.Stop(); err != nil {
				s.log.Errorf("Failed to stop user lifecycle: %v", err)
				return err
			}
			s.log.Infof("user lifecycle stopped.")
			return nil
		},
	})
}
