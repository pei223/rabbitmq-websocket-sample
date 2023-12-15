package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	Logger       zerolog.Logger
	AccessLogger zerolog.Logger
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	Logger = zerolog.New(
		os.Stdout,
	).With().Timestamp().Caller().Logger()
	AccessLogger = zerolog.New(
		os.Stdout,
	).With().Timestamp().Caller().Logger()
}

func SetLogLevel(logLevelStr string) {
	logLevel, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		Logger.Error().Msgf("Invalid log level: %s. Set info level.", logLevelStr)
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
}
