package logging

import (
	"os"

	"github.com/rs/zerolog"
)

func New(isDebug bool) *zerolog.Logger {
	logLvl := zerolog.InfoLevel
	if isDebug {
		logLvl = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(logLvl)
	// output := zerolog.ConsoleWriter{
	// 	// Out:        ,
	// 	TimeFormat: time.RFC3339,
	// 	NoColor:    true,
	// }
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &logger
}
