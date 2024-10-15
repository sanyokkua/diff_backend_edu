package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestGetEnvWithDefault_ValueSet(t *testing.T) {
	_ = os.Setenv("TEST_ENV", "value")
	defer os.Unsetenv("TEST_ENV")

	result := getEnvWithDefault("TEST_ENV", "default")
	assert.Equal(t, "value", result)
}

func TestGetEnvWithDefault_ValueNotSet(t *testing.T) {
	result := getEnvWithDefault("TEST_ENV_NOT_SET", "default")
	assert.Equal(t, "default", result)
}

func TestGetConfig(t *testing.T) {
	_ = os.Setenv(DbHost, "localhost")
	_ = os.Setenv(DbUser, "user")
	_ = os.Setenv(DbPassword, "password")
	_ = os.Setenv(DbName, "dbname")
	_ = os.Setenv(DbPort, "5432")
	_ = os.Setenv(JwtSecret, "secret")
	defer os.Clearenv()

	config := GetConfig()

	assert.Equal(t, "localhost", config[DbHost])
	assert.Equal(t, "user", config[DbUser])
	assert.Equal(t, "password", config[DbPassword])
	assert.Equal(t, "dbname", config[DbName])
	assert.Equal(t, "5432", config[DbPort])
	assert.Equal(t, "secret", config[JwtSecret])
}

func TestGetDbDsn(t *testing.T) {
	config := map[string]string{
		DbHost:     "localhost",
		DbUser:     "user",
		DbPassword: "password",
		DbName:     "dbname",
		DbPort:     "5432",
	}

	dsn := GetDbDsn(config)
	expected := "host=localhost user=user password=password dbname=dbname port=5432 sslmode=disable"
	assert.Equal(t, expected, dsn)
}

func TestGetJwtSecret(t *testing.T) {
	config := map[string]string{
		JwtSecret: "supersecret",
	}

	secret := GetJwtSecret(config)
	assert.Equal(t, "supersecret", secret)
}

func TestGetCorsConfig(t *testing.T) {
	corsConfig := GetCorsConfig()

	assert.Contains(t, corsConfig.AllowOrigins, "http://localhost:5173")
	assert.Contains(t, corsConfig.AllowMethods, "GET")
	assert.Equal(t, corsConfig.AllowCredentials, true)
	assert.Equal(t, corsConfig.MaxAge, 12*time.Hour)
}
