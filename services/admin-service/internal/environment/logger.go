package environment

import (
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	appLogger   *slog.Logger
	loggerOnce  sync.Once
	serviceName string
	environment string
)

func InitLogger(env, svcName string) {
	loggerOnce.Do(func() {
		serviceName = svcName
		environment = env

		var handler slog.Handler
		handlerOptions := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey && len(groups) == 0 {
					return slog.Attr{Key: "timestamp", Value: slog.StringValue(a.Value.Any().(time.Time).Format(time.RFC3339))}
				}
				if a.Key == slog.LevelKey && len(groups) == 0 {
					return slog.Attr{Key: "log_level", Value: a.Value}
				}
				if a.Key == slog.MessageKey && len(groups) == 0 {
					return slog.Attr{Key: "message", Value: a.Value}
				}
				return a
			},
		}

		switch env {
		case "development":
			handler = slog.NewTextHandler(os.Stdout, handlerOptions)
		case "production":
			handler = slog.NewJSONHandler(os.Stderr, handlerOptions)
			handlerOptions.Level = slog.LevelInfo
		default:
			slog.Warn("Unknown environment, defaulting to development logging level")
			handler = slog.NewTextHandler(os.Stdout, handlerOptions)
		}

		appLogger = slog.New(handler).With(
			slog.String("service", serviceName),
			slog.String("environment", environment))

	})
}

func GetLogger() *slog.Logger {
	if appLogger == nil {
		slog.Error("Logger not initialized, call InitLogger first")
		return nil
	}
	return appLogger
}
