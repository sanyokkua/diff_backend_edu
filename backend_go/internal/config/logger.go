package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

// InitLogger initializes the global logger with default settings.
// It configures the logger to output to standard output with timestamps, caller information,
// and sets the global log level to Info.
//
// This function is typically called during the application startup to ensure that logging
// is properly set up before any log messages are generated.
//
// It uses the zerolog package for structured logging and formats the time field in RFC3339 format.
func InitLogger() {
	// Set global logger with timestamp and human-readable output for development
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Caller().
		Logger()

	// Set the global log level to info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Optionally, set time format for logs
	zerolog.TimeFieldFormat = time.RFC3339
}
