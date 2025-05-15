package order

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

// ServerStateModule defines the Fx module for the order server
var ServerStateModule = fx.Module("order",
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
	fiberApp := provider.NewFiberApp(v, log, "order")

	fiberApp.Get("/order/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "order server running"})
	})
	return &ServerState{
		log:      log,
		fiberApp: fiberApp,
		v:        v,
	}
}

func (s *ServerState) Start() error {
	pwd, _ := os.Getwd()
	cert := pwd + s.v.GetString("order.certificate.cert")
	key := pwd + s.v.GetString("order.certificate.key")
	port := os.Getenv("ORDER_PORT")

	s.log.Infof("Starting order server on port %s...", port)

	go func() {
		if s.v.GetString("app.env") == "production" {
			if err := s.fiberApp.ListenTLS(":"+port, cert, key); err != nil {
				s.log.Errorf("Failed to start order server: %v", err)
			}
		} else {
			if err := s.fiberApp.Listen(":" + port); err != nil {
				s.log.Errorf("Failed to start order server: %v", err)
			}
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

// RegisterLifeCycle registers the lifecycle with Fx
func RegisterLifeCycle(lc fx.Lifecycle, s *ServerState) {
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
