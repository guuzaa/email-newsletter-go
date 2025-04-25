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
}
