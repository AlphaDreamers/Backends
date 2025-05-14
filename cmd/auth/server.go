package auth

import (
	"context"
	"github.com/SwanHtetAungPhyo/backend/internal/handler/auth"
	"github.com/SwanHtetAungPhyo/backend/test_mulit/provider"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var ServerStateModule = fx.Module("auth",
	fx.Provide(NewServerState),
	fx.Invoke(RegisterLifeCycle),
)

type ServerState struct {
	log      *logrus.Logger
	fiberApp *fiber.App
	v        *viper.Viper
	handler  *auth.Handler
}

func NewServerState(
	log *logrus.Logger,
	v *viper.Viper,
	handler *auth.Handler,
) *ServerState {
	fiberApp := provider.NewFiberApp(v, log, "auth")
	fiberApp.Get("/auth/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "auth server running"})
	})
	return &ServerState{
		log:      log,
		fiberApp: fiberApp,
		v:        v,
		handler:  handler,
	}
}

func (s *ServerState) Start() error {
	pwd, _ := os.Getwd()
	cert := pwd + s.v.GetString("auth.certificate.cert")
	key := pwd + s.v.GetString("auth.certificate.key")
	port := s.v.GetString("auth.port")

	s.log.Infof("Starting auth server on port %s...", port)

	go func() {
		if s.v.GetString("app.env") == "production" {
			if err := s.fiberApp.ListenTLS(":"+port, cert, key); err != nil {
				s.log.Errorf("Failed to start auth server: %v", err)
			}
		} else {
			if err := s.fiberApp.Listen(":" + port); err != nil {
				s.log.Errorf("Failed to start auth server: %v", err)
			}
		}
	}()
	return nil
}

func (s *ServerState) Stop() error {
	s.log.Infof("Shutting down auth server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		s.log.Errorf("Failed to gracefully shut down auth server: %v", err)
		return err
	}

	s.log.Infof("auth server shutdown complete.")
	return nil
}

// RegisterLifeCycle registers the lifecycle with Fx
func RegisterLifeCycle(lc fx.Lifecycle, s *ServerState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Infof("Starting auth lifecycle...")
			if err := s.Start(); err != nil {
				s.log.Errorf("Failed to start auth lifecycle: %v", err)
				return err
			}
			s.log.Infof("auth lifecycle started.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Infof("Stopping auth lifecycle...")
			if err := s.Stop(); err != nil {
				s.log.Errorf("Failed to stop auth lifecycle: %v", err)
				return err
			}
			s.log.Infof("auth lifecycle stopped.")
			return nil
		},
	})
}
