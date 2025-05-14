package wallet

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var ServerStateModule = fx.Module("wallet", fx.Provide(
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

	s.log.Infof("Starting wallet server on port %s...", port)

	go func() {
		if err := s.fiberApp.ListenTLS(":" + port, cert, key); err != nil {
			s.log.Errorf("Failed to start server: %v", err)
		}
	}()
	return nil
}

func (s *ServerState) Stop() error {
	s.log.Infof("Shutting down wallet server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		s.log.Errorf("Failed to gracefully shut down wallet server: %v", err)
		return err
	}

	s.log.Infof("wallet server shutdown complete.")
	return nil
}

// Register lifecycle with fx using fx.Hook
func RegisterLifewalletCycle(lc fx.Lifecycle, s *ServerState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Infof("Starting wallet lifecycle...")
			if err := s.Start(); err != nil {
				s.log.Errorf("Failed to start wallet lifecycle: %v", err)
				return err
			}
			s.log.Infof("wallet lifecycle started.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Infof("Stopping wallet lifecycle...")
			if err := s.Stop(); err != nil {
				s.log.Errorf("Failed to stop wallet lifecycle: %v", err)
				return err
			}
			s.log.Infof("wallet lifecycle stopped.")
			return nil
		},
	})
}
