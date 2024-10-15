package config

import (
	"bytes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

// TestInitLogger verifies that InitLogger correctly configures the global logger.
func TestInitLogger(t *testing.T) {
	// Capture the current stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call InitLogger to set up the logger
	InitLogger()

	// Verify the global log level is set to Info
	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())

	// Create a logger message to test if output is written to stdout
	log.Info().Msg("Test log message")

	// Flush the pipe to get the log output
	_ = w.Close()
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	assert.NoError(t, err)

	// Reset stdout
	os.Stdout = oldStdout

	// Check that the log output contains the expected elements: "Test log message"
	assert.Contains(t, buf.String(), "Test log message")

	// Check that the log contains a timestamp in RFC3339 format
	expectedTimeFormat := time.Now().Format(time.RFC3339)[:10] // Checking only the date part
	assert.Contains(t, buf.String(), expectedTimeFormat)

	// Check that the log contains caller information (file path and line number)
	assert.Contains(t, buf.String(), "logger_test.go")
}
