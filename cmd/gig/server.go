package gig

import (
	"context"
	"github.com/SwanHtetAungPhyo/backend/cmd/provider"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// ServerStateModule defines the Fx module for the gig server
var ServerStateModule = fx.Module("gig",
	fx.Provide(
		NewServerState,
	),
	fx.Invoke(
		RegisterLifeCycle,
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
	fiberApp := provider.NewFiberApp(v, log, "gig")
	// Define a basic endpoint for the gig server
	fiberApp.Get("/gig/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "gig server running"})
	})
	return &ServerState{
		log:      log,
		fiberApp: fiberApp,
		v:        v,
	}
}

func (s *ServerState) Start() error {
	pwd, _ := os.Getwd()
	cert := pwd + s.v.GetString("gig.certificates.cert")
	key := pwd + s.v.GetString("gig.certificates.key")
	port := os.Getenv("GIG_PORT")

	s.log.Infof("Starting gig server on port %s...", port)

	go func() {
		if s.v.GetString("app.env") == "production" {
			if err := s.fiberApp.ListenTLS(":"+port, cert, key); err != nil {
				s.log.Errorf("Failed to start gig server: %v", err)
			}
		} else {
			if err := s.fiberApp.Listen(":" + port); err != nil {
				s.log.Errorf("Failed to start gig server: %v", err)
			}
		}
	}()
	return nil
}

func (s *ServerState) Stop() error {
	s.log.Infof("Shutting down gig server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		s.log.Errorf("Failed to gracefully shut down gig server: %v", err)
		return err
	}

	s.log.Infof("gig server shutdown complete.")
	return nil
}

// RegisterLifeCycle registers the lifecycle with Fx
func RegisterLifeCycle(lc fx.Lifecycle, s *ServerState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Infof("Starting gig lifecycle...")
			if err := s.Start(); err != nil {
				s.log.Errorf("Failed to start gig lifecycle: %v", err)
				return err
			}
			s.log.Infof("gig lifecycle started.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Infof("Stopping gig lifecycle...")
			if err := s.Stop(); err != nil {
				s.log.Errorf("Failed to stop gig lifecycle: %v", err)
				return err
			}
			s.log.Infof("gig lifecycle stopped.")
			return nil
		},
	})
}
