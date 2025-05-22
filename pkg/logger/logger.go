package logger

import (
	"fmt"
	"github.com/kkboranbay/task-service/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

func SetupLogger(cfg config.LoggerConfig) {
	zerolog.TimeFieldFormat = time.RFC3339

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| %s |", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
}

// глобальный логгер для удобства
func L() *zerolog.Logger {
	return &log.Logger
}

// логгер с добавленным полем
func WithField(key string, value interface{}) *zerolog.Logger {
	logger := log.With().Interface(key, value).Logger()
	return &logger
}

// логгер с добавленными полями
func WithFields(fields map[string]interface{}) *zerolog.Logger {
	ctx := log.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	logger := ctx.Logger()
	return &logger
}

// логгер с добавленной ошибкой
func WithError(err error) *zerolog.Logger {
	logger := log.With().Err(err).Logger()
	return &logger
}
