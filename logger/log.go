package logger

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	loggerWithoutCaller = zerolog.Nop()
	once                sync.Once
)

func Init(debug bool, pretty bool, additionalWriters ...io.Writer) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if pretty {
		additionalWriters = append(additionalWriters, zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		additionalWriters = append(additionalWriters, os.Stderr)
	}

	once.Do(func() {
		loggerWithoutCaller = log.Output(zerolog.MultiLevelWriter(additionalWriters...))
	})
	log.Logger = loggerWithoutCaller.With().Caller().Logger()
}

func WithoutCaller() zerolog.Logger {
	return loggerWithoutCaller
}
