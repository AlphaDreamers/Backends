// provider/fiber.go
package provider

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	flog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

func SetLogger(v *viper.Viper) *logrus.Logger {
	logger := logrus.New()
	logFile := v.GetString("log.file")
	if logFile == "" {
		logDir := v.GetString("log.dir")
		if logDir == "" {
			logDir = "./logs"
		}
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
		}
		logFile = path.Join(logDir, "app.log")
	}
	lumberJack := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    v.GetInt("log.max_size"), // MB
		MaxBackups: v.GetInt("log.max_backups"),
		MaxAge:     v.GetInt("log.max_age"), // days
		Compress:   v.GetBool("log.compress"),
	}
	multiWriter := io.MultiWriter(os.Stdout, lumberJack)
	logger.SetOutput(multiWriter)
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		PrettyPrint: v.GetBool("log.pretty_print"),
	})
	logger.SetReportCaller(true)
	return logger
}

// NewFiberApp creates a production-ready Fiber app
func NewFiberApp(v *viper.Viper, logger *logrus.Logger) *fiber.App {
	isProduction := v.GetString("app.env") == "production"

	// Fiber configuration
	app := fiber.New(fiber.Config{
		DisableStartupMessage: v.GetBool("fiber.disableStartupMessage"),
		StrictRouting:         v.GetBool("fiber.strictRouting"),
		CaseSensitive:         v.GetBool("fiber.caseSensitive"),
		ServerHeader:          v.GetString("fiber.serverHeader"),
		AppName:               v.GetString("fiber.appName"),
		ReduceMemoryUsage:     true,
		Prefork:               v.GetBool("fiber.prefork") && isProduction,
		ReadTimeout:           v.GetDuration("fiber.readTimeout"),
		WriteTimeout:          v.GetDuration("fiber.writeTimeout"),
		IdleTimeout:           v.GetDuration("fiber.idleTimeout"),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logger.WithFields(logrus.Fields{
				"path": ctx.Path(),
			}).Errorf("Request failed: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware: Recover from panics
	app.Use(recover.New())

	// Middleware: Security headers
	app.Use(helmet.New())

	// Middleware: CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     v.GetString("fiber.allowed-origin"),
		AllowHeaders:     v.GetString("fiber.allowHeaders"),
		AllowMethods:     v.GetString("fiber.allowMethods"),
		AllowCredentials: v.GetBool("fiber.allowCredentials"),
		MaxAge:           86400,
	}))

	// Middleware: CSRF protection (skip in development)
	if isProduction {
		app.Use(csrf.New(csrf.Config{
			KeyLookup:      "header:X-Csrf-Token",
			CookieName:     "csrf_",
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: "Strict",
			Expiration:     1 * time.Hour,
			ContextKey:     "csrf",
			Extractor:      csrf.CsrfFromHeader("X-Csrf-Token"),
			KeyGenerator:   func() string { return uuid.New().String() },
		}))
	}

	app.Use(limiter.New(limiter.Config{
		Max:        v.GetInt("limiter.max_requests"),
		Expiration: v.GetDuration("limiter.expiration"),
		LimitReached: func(c *fiber.Ctx) error {
			logger.WithFields(logrus.Fields{
				"ip":   c.IP(),
				"path": c.Path(),
			}).Warn("Rate limit exceeded")
			return fiber.ErrTooManyRequests
		},
		SkipFailedRequests: true,
	}))

	// Middleware: Request logging
	app.Use(flog.New(flog.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "UTC",
		Output:     &logrusWriter{logger: logger},
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	return app
}

type logrusWriter struct {
	logger *logrus.Logger
}

func (w *logrusWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}
