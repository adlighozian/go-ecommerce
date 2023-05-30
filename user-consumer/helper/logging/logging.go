package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(isDebug bool) *zerolog.Logger {
	logLvl := zerolog.InfoLevel
	if isDebug {
		logLvl = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(logLvl)
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &logger
}
