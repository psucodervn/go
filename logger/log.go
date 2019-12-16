package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	log.Logger = log.Output(zerolog.MultiLevelWriter(additionalWriters...)).With().Caller().Logger()
}
