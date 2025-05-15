// provider/fiber.go
package provider

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	flog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

// SetLogger configures a logrus logger with file rotation
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
		MaxSize:    v.GetInt("log.max_size"),
		MaxBackups: v.GetInt("log.max_backups"),
		MaxAge:     v.GetInt("log.max_age"),
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

// NewFiberApp creates a Fiber app with basic configuration
func NewFiberApp(v *viper.Viper, logger *logrus.Logger, prefix string) *fiber.App {
	isProduction := v.GetString("app.env") == "production"
	//pwd, err := os.Getwd()
	//if err != nil {
	//	logger.Fatal(err.Error())
	//}
	//viewPath := pwd + "/cmd" + "/views"
	//logrus.Info("View path: " + viewPath)
	//views := html.New(viewPath, ".html")
	app := fiber.New(fiber.Config{
		//Views:                 views,
		DisableStartupMessage: isProduction,
		ServerHeader:          v.GetString(prefix + ".header"),
		AppName:               v.GetString(prefix + ".name"),
		ReduceMemoryUsage:     true,
		ReadTimeout:           v.GetDuration(prefix + ".read_timeout"),
		WriteTimeout:          v.GetDuration(prefix + ".write_timeout"),
		IdleTimeout:           v.GetDuration(prefix + ".idle_timeout"),
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

	// Middleware: Request logging
	app.Use(flog.New(flog.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "UTC",
		Output:     &logrusWriter{logger: logger},
	}))

	return app
}

// logrusWriter adapts logrus for Fiber's logger middleware
type logrusWriter struct {
	logger *logrus.Logger
}

func (w *logrusWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}
