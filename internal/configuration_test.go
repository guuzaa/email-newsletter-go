package internal_test

import (
	"os"
	"testing"

	"github.com/guuzaa/email-newsletter/internal"
	"github.com/stretchr/testify/assert"
)

func TestConfiguration(t *testing.T) {
	os.Setenv("APP_ENVIRONMENT", "local")
	os.Setenv("APP_HOST", "127.0.0.2")
	settings, err := internal.Configuration("../configuration")
	assert.Nil(t, err, "Failed to load configuration")
	assert.Equal(t, "postgres", settings.Database.Username)
	assert.Equal(t, "password", settings.Database.Password)
	assert.Equal(t, uint16(5432), settings.Database.Port)
	assert.Equal(t, "127.0.0.1", settings.Database.Host)
	assert.Equal(t, "newsletter", settings.Database.DatabaseName)
	assert.Equal(t, "127.0.0.2", settings.Application.Host)
	assert.Equal(t, uint16(8000), settings.Application.Port)
	assert.False(t, settings.Database.RequireSSL)
	t.Cleanup(func() {
		os.Unsetenv("APP_ENVIRONMENT")
		os.Unsetenv("APP_HOST")
	})
}

func TestConfigurationWithMissingFile(t *testing.T) {
	os.Setenv("APP_ENVIRONMENT", "local")
	os.Setenv("APP_HOST", "127.0.0.2")
	os.Setenv("APP_PORT", "9000")
	os.Setenv("APP_DB_USERNAME", "test")
	os.Setenv("APP_DB_PASSWORD", "test")
	os.Setenv("APP_DB_PORT", "54327")
	os.Setenv("APP_DB_HOST", "127.0.0.2")
	os.Setenv("APP_DB_NAME", "newsletter")
	os.Setenv("APP_DB_REQUIRE_SSL", "false")
	settings, err := internal.Configuration("404notfound404")
	assert.Nil(t, err, "Failed to load configuration")
	assert.Equal(t, "test", settings.Database.Username)
	assert.Equal(t, "test", settings.Database.Password)
	assert.Equal(t, uint16(54327), settings.Database.Port)
	assert.Equal(t, "127.0.0.2", settings.Database.Host)
	assert.Equal(t, "newsletter", settings.Database.DatabaseName)
	assert.False(t, settings.Database.RequireSSL)
	t.Cleanup(func() {
		os.Unsetenv("APP_ENVIRONMENT")
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_DATABASE_NAME")
		os.Unsetenv("DB_REQUIRE_SSL")
	})
}

func TestConfigurationWithMissingEnvironmentVariables(t *testing.T) {
	_, err := internal.Configuration("404notfound404")
	assert.NotNil(t, err, "Failed to load configuration")
}
