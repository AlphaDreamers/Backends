package auth

import (
	"context"
	"github.com/SwanHtetAungPhyo/backend/internal/handler/auth"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var ServerStateModule = fx.Module("auth", fx.Provide(
	NewServerState,
))

type ServerState struct {
	log      *logrus.Logger
	fiberApp *fiber.App
	v        *viper.Viper
	handler  *auth.Handler
}

func NewServerState(
	log *logrus.Logger,
	fiberApp *fiber.App,
	v *viper.Viper,
	handler *auth.Handler,
) *ServerState {
	return &ServerState{
		log:      log,
		fiberApp: fiberApp,
		v:        v,
		handler:  handler,
	}
}

func (s *ServerState) Start() error {
	pwd, _ := os.Getwd()
	cert := pwd + s.v.GetString("fiber.certificate.cert")
	key := pwd + s.v.GetString("fiber.certificate.key")
	port := s.v.GetString("fiber.auth.port")

	s.log.Infof("Starting auth server on port %s...", port)

	go func() {
		if err := s.fiberApp.ListenTLS(":"+port, cert, key); err != nil {
			s.log.Errorf("Failed to start server: %v", err)
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

// RegisterLifeboatCycle Register lifecycle with fx using fx.Hook
func RegisterLifeauthCycle(lc fx.Lifecycle, s *ServerState) {
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
