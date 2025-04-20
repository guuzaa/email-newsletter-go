package internal_test

import (
	"testing"

	"github.com/guuzaa/email-newsletter/internal"
	"github.com/stretchr/testify/assert"
)

func TestConfiguration(t *testing.T) {
	settings, err := internal.Configuration("../configuration.yaml")
	assert.Nil(t, err, "Failed to load configuration")
	assert.Equal(t, "postgres", settings.Username)
	assert.Equal(t, "password", settings.Password)
	assert.Equal(t, uint16(5432), settings.Port)
	assert.Equal(t, "127.0.0.1", settings.Host)
	assert.Equal(t, "newsletter", settings.DatabaseName)
}
