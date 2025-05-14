package order

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var ServerStateModule = fx.Module("order", fx.Provide(
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

	s.log.Infof("Starting order server on port %s...", port)

	go func() {
		if err := s.fiberApp.ListenTLS(":" + port, cert, key); err != nil {
			s.log.Errorf("Failed to start server: %v", err)
		}
	}()
	return nil
}

func (s *ServerState) Stop() error {
	s.log.Infof("Shutting down order server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		s.log.Errorf("Failed to gracefully shut down order server: %v", err)
		return err
	}

	s.log.Infof("order server shutdown complete.")
	return nil
}

// Register lifecycle with fx using fx.Hook
func RegisterLifeorderCycle(lc fx.Lifecycle, s *ServerState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Infof("Starting order lifecycle...")
			if err := s.Start(); err != nil {
				s.log.Errorf("Failed to start order lifecycle: %v", err)
				return err
			}
			s.log.Infof("order lifecycle started.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Infof("Stopping order lifecycle...")
			if err := s.Stop(); err != nil {
				s.log.Errorf("Failed to stop order lifecycle: %v", err)
				return err
			}
			s.log.Infof("order lifecycle stopped.")
			return nil
		},
	})
}
