package provider

//
//import (
//	"fmt"
//	"github.com/sirupsen/logrus"
//	"gopkg.in/natefinch/lumberjack.v2"
//	"io"
//	"os"
//	"path"
//	"runtime"
//)
//
//func SetLogger() *logrus.Logger {
//	logger := logrus.New()
//	lumperJack := &lumberjack.Logger{
//		Filename:   fmt.Sprintf("%s/%s.log", os.Getenv("HOME"), os.Getenv("LOG")),
//		MaxSize:    10,
//		MaxBackups: 3,
//		MaxAge:     28,
//		Compress:   true,
//	}
//	multiWriter := io.MultiWriter(os.Stdout, lumperJack)
//	logger.SetOutput(multiWriter)
//	logger.SetFormatter(&logrus.JSONFormatter{
//		FieldMap: logrus.FieldMap{
//			logrus.FieldKeyTime:  "timestamp",
//			logrus.FieldKeyLevel: "severity",
//			logrus.FieldKeyMsg:   "message",
//		},
//		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
//			filename := path.Base(f.File)
//			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
//		},
//		PrettyPrint: true,
//	})
//	return logger
//}
