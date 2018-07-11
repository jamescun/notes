package main

import (
	"net/http"
	"os"

	"github.com/jamescun/notes"

	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Timestamp().
		Logger()

	cfg, err := notes.ConfigFromEnvironment()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config from environment")
	}

	s := &http.Server{
		Addr:    cfg.Addr,
		Handler: notes.Handler(log, cfg),
	}

	log.Info().Str("addr", s.Addr).Msg("server listening")

	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
