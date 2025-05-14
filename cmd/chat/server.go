package chat

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var ServerStateModule = fx.Module("chat", fx.Provide(
	NewServerState,
))

type ServerState struct {
	log      *logrus.Logger
	fiberApp *fiber.App
	v        *viper.Viper
}

func NewServerState(
	log *logrus.Logger,
	fiberApp *fiber.App,
	v *viper.Viper,
) *ServerState {
	return &ServerState{
		log:      log,
		fiberApp: fiberApp,
		v:        v,
	}
}

func (s *ServerState) Start() error {
	pwd, _ := os.Getwd()
	cert := pwd + s.v.GetString("fiber.certificate.cert")
	key := pwd + s.v.GetString("fiber.certificate.key")
	port := s.v.GetString("fiber.port")

	s.log.Infof("Starting chat server on port %s...", port)

	go func() {
		if err := s.fiberApp.ListenTLS(":" + port, cert, key); err != nil {
			s.log.Errorf("Failed to start server: %v", err)
		}
	}()
	return nil
}

func (s *ServerState) Stop() error {
	s.log.Infof("Shutting down chat server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		s.log.Errorf("Failed to gracefully shut down chat server: %v", err)
		return err
	}

	s.log.Infof("chat server shutdown complete.")
	return nil
}

// Register lifecycle with fx using fx.Hook
func RegisterLifechatCycle(lc fx.Lifecycle, s *ServerState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Infof("Starting chat lifecycle...")
			if err := s.Start(); err != nil {
				s.log.Errorf("Failed to start chat lifecycle: %v", err)
				return err
			}
			s.log.Infof("chat lifecycle started.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Infof("Stopping chat lifecycle...")
			if err := s.Stop(); err != nil {
				s.log.Errorf("Failed to stop chat lifecycle: %v", err)
				return err
			}
			s.log.Infof("chat lifecycle stopped.")
			return nil
		},
	})
}
