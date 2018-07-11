package notes

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/jamescun/notes/db"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

// Config contains options for configuring the operation of the Notes app.
type Config struct {
	Addr string `json:"addr"`

	DSN string `json:"dsn"`
}

// DefaultConfig returns the Notes app config seeded with default values.
func DefaultConfig() Config {
	return Config{
		Addr: "127.0.0.1:8080",
		DSN:  "postgres://postgres:lolsecurity@localhost/notes?sslmode=disable",
	}
}

// ConfigFromEnvironment unmarshals the Notes app config from an environment
// variable called 'CONFIG', missing values will be seeded with defaults.
func ConfigFromEnvironment() (Config, error) {
	cfg := DefaultConfig()

	env := os.Getenv("CONFIG")
	if env == "" {
		return cfg, nil
	}

	err := json.Unmarshal([]byte(os.Getenv("CONFIG")), &cfg)
	return cfg, err
}

// Handler returns a configured http Handler for the Notes app.
func Handler(log zerolog.Logger, cfg Config) http.Handler {
	log = log.With().Str("service", "notes").Logger()

	dbConn, err := db.New(cfg.DSN)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to database")
	}
	defer dbConn.Close()

	mux := chi.NewRouter()

	return mux
}
