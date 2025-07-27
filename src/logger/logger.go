package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	globalLogger zerolog.Logger
)

func getLogLevel() zerolog.Level {
	logLevelStr := os.Getenv("LOG_LEVEL")
	appEnv := strings.ToLower(os.Getenv("APP_ENV"))

	if logLevelStr != "" {
		if level, err := zerolog.ParseLevel(logLevelStr); err == nil {
			return level
		}
	}

	switch appEnv {
	case "production", "prod":
		return zerolog.WarnLevel
	case "uat", "sit":
		return zerolog.InfoLevel
	default:
		return zerolog.DebugLevel
	}
}

func Init(serviceName, env, logLevelStr string) {
	zerolog.TimeFieldFormat = time.RFC3339
	level, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		// fallback jika parsing gagal, bisa set default berdasarkan env
		switch strings.ToLower(env) {
		case "production", "prod":
			level = zerolog.WarnLevel
		case "uat", "sit":
			level = zerolog.InfoLevel
		default:
			level = zerolog.DebugLevel
		}
	}
	zerolog.SetGlobalLevel(level)
	var output = os.Stdout
	var baseLogger zerolog.Logger

	if env == "dev" {
		baseLogger = zerolog.New(zerolog.ConsoleWriter{Out: output}).
			With().
			Timestamp().
			Str("service", serviceName).
			Logger()
	} else {
		baseLogger = zerolog.New(output).
			With().
			Timestamp().
			Str("service", serviceName).
			Logger()
	}

	globalLogger = baseLogger
	log.Logger = globalLogger
}

// Get returns the base global logger
func Get() zerolog.Logger {
	return globalLogger
}

// WithContext attaches logger with request context fields
func WithContext(ctx context.Context) *zerolog.Logger {
	logger := zerolog.Ctx(ctx)
	if logger == nil {
		return &log.Logger
	}
	return logger
}

func InjectContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}
