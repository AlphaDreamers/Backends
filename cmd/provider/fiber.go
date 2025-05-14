package provider

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"time"
)

func NewFiberApp(v *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		StrictRouting:         false,
		ServerHeader:          v.GetString("server.header"),
		AppName:               v.GetString("server.name"),
		ReduceMemoryUsage:     true,
		ReadTimeout:           v.GetDuration("server.read_timeout"),
		WriteTimeout:          v.GetDuration("server.write_timeout"),
		IdleTimeout:           v.GetDuration("server.idle_timeout"),
	})
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return true
		},
		Max:        5,
		Expiration: 5 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	}))
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     v.GetString("app.allow_origins"),
		AllowHeaders:     v.GetString("app.allow_headers"),
		AllowMethods:     v.GetString("app.allow_methods"),
		AllowCredentials: v.GetBool("app.allow_credentials"),
	}))
	return app
}
