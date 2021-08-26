package internal

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Ddd(is string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Print(is)

}
