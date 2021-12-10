package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func initLog(verbose bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:          os.Stderr,
		PartsExclude: []string{"time", "level"},
	})
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
