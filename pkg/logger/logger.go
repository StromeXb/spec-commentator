package logger

import (
	"os"
	"time"

	"github.com/spf13/viper"

	"spec-commentor/pkg/test"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EnvLogLevel = "LOG_LEVEL"
)

type Config struct {
	LogLevel string
}

func NewLoggerConfig(v *viper.Viper) *Config {
	return &Config{
		LogLevel: v.GetString(EnvLogLevel),
	}
}

func NewLogger(cfg *Config) (*zerolog.Logger, error) {
	lvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed parse log level")

		return nil, err
	}
	logger := zerolog.New(os.Stderr).Level(lvl).With().Caller().Timestamp().Logger()
	if !test.IsTest() {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).Level(lvl).With().Caller().Timestamp().Logger()
	}
	return &logger, nil
}
