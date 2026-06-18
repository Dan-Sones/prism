package prismLog

import (
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	appLogger   *slog.Logger
	loggerOnce  sync.Once
	serviceName string
	environment string
)

func InitLogger(env, svcName, level, format string) {
	loggerOnce.Do(func() {
		serviceName = svcName
		environment = env

		handlerOptions := &slog.HandlerOptions{
			AddSource: true,
			Level:     parseLevel(level),
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

		var handler slog.Handler
		switch strings.ToLower(format) {
		case "json":
			handler = slog.NewJSONHandler(os.Stderr, handlerOptions)
		case "text":
			handler = slog.NewTextHandler(os.Stdout, handlerOptions)
		default:
			slog.Warn("Unknown log format, defaulting to text", slog.String("format", format))
			handler = slog.NewTextHandler(os.Stdout, handlerOptions)
		}

		appLogger = slog.New(handler).With(
			slog.String("service", serviceName),
			slog.String("environment", environment))
	})
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Warn("Unknown log level, defaulting to info", slog.String("level", level))
		return slog.LevelInfo
	}
}

func GetLogger() *slog.Logger {
	if appLogger == nil {
		slog.Error("Logger not initialized, call InitLogger first")
		return nil
	}
	return appLogger
}
