package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/guuzaa/email-newsletter/internal/domain"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	Database    DatabaseSettings    `yaml:"database"`
	Application ApplicationSettings `yaml:"application"`
	EmailClient EmailClientSettings `yaml:"email_client"`
}

type ApplicationSettings struct {
	Port uint16 `yaml:"port" env:"APP_PORT"`
	Host string `yaml:"host" env:"APP_HOST"`
}

type EmailClientSettings struct {
	BaseURL             string `yaml:"base_url" env:"APP_EMAIL_BASE_URL"`
	SenderEmail         string `yaml:"sender_email" env:"APP_SENDER_EMAIL"`
	AuthorizationToken  string `yaml:"authorization_token" env:"APP_EMAIL_AUTHORIZATION_TOKEN"`
	TimeoutMilliseconds uint64 `yaml:"timeout_milliseconds" env:"APP_EMAIL_CLIENT_TIMEOUT_MILLISECONDS"`
}

func (ecs EmailClientSettings) Sender() (domain.SubscriberEmail, error) {
	return domain.SubscriberEmailFrom(ecs.SenderEmail)
}

func (ecs EmailClientSettings) Timeout() time.Duration {
	return time.Duration(ecs.TimeoutMilliseconds) * time.Millisecond
}

type DatabaseSettings struct {
	Username     string `yaml:"username" env:"APP_DB_USERNAME"`
	Password     string `yaml:"password" env:"APP_DB_PASSWORD"`
	Port         uint16 `yaml:"port" env:"APP_DB_PORT"`
	Host         string `yaml:"host" env:"APP_DB_HOST"`
	DatabaseName string `yaml:"database_name" env:"APP_DB_NAME"`
	RequireSSL   bool   `yaml:"require_ssl" env:"APP_DB_REQUIRE_SSL"`
}

func (setting Settings) PostgresSQLDSN() string {
	sslMode := "disable"
	if setting.Database.RequireSSL {
		sslMode = "require"
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", setting.Database.Host, setting.Database.Username, setting.Database.Password, setting.Database.DatabaseName, setting.Database.Port, sslMode)
}

func (setting *Settings) Valid() bool {
	return setting.Application.Host != "" && setting.Application.Port != 0 &&
		setting.Database.Host != "" && setting.Database.Port != 0 &&
		setting.Database.Username != "" && setting.Database.Password != "" &&
		setting.Database.DatabaseName != ""
}

func (setting Settings) Address() string {
	return fmt.Sprintf("%s:%d", setting.Application.Host, setting.Application.Port)
}

// mergeSettings combines the settings from the target into the base settings
// only if they are not empty or zero values
func mergeSettings(base, overlay Settings) Settings {
	result := base

	if overlay.Database.Host != "" {
		result.Database.Host = overlay.Database.Host
	}
	if overlay.Database.Port != 0 {
		result.Database.Port = overlay.Database.Port
	}
	if overlay.Database.Username != "" {
		result.Database.Username = overlay.Database.Username
	}
	if overlay.Database.Password != "" {
		result.Database.Password = overlay.Database.Password
	}
	if overlay.Database.DatabaseName != "" {
		result.Database.DatabaseName = overlay.Database.DatabaseName
	}
	if overlay.Database.RequireSSL {
		result.Database.RequireSSL = overlay.Database.RequireSSL
	}

	if overlay.Application.Port != 0 {
		result.Application.Port = overlay.Application.Port
	}
	if overlay.Application.Host != "" {
		result.Application.Host = overlay.Application.Host
	}

	if overlay.EmailClient.BaseURL != "" {
		result.EmailClient.BaseURL = overlay.EmailClient.BaseURL
	}
	if overlay.EmailClient.SenderEmail != "" {
		result.EmailClient.SenderEmail = overlay.EmailClient.SenderEmail
	}
	if overlay.EmailClient.AuthorizationToken != "" {
		result.EmailClient.AuthorizationToken = overlay.EmailClient.AuthorizationToken
	}
	if overlay.EmailClient.TimeoutMilliseconds != 0 {
		result.EmailClient.TimeoutMilliseconds = overlay.EmailClient.TimeoutMilliseconds
	}

	return result
}

func Configuration(path string) (Settings, error) {
	var settings Settings
	logger := Logger()

	baseFilePath := filepath.Join(path, "base.yaml")
	data, err := os.ReadFile(baseFilePath)
	if err == nil {
		yaml.Unmarshal(data, &settings)
	}

	environment := ParseEnvironment(os.Getenv("APP_ENVIRONMENT"))
	envFilePath := filepath.Join(path, fmt.Sprintf("%s.yaml", environment.String()))
	data, err = os.ReadFile(envFilePath)
	var envSettings Settings
	if err == nil {
		yaml.Unmarshal(data, &envSettings)
	}
	settings = mergeSettings(settings, envSettings)

	var envVarSettings Settings
	if err := env.Parse(&envVarSettings); err != nil {
		logger.Debug().Err(err).Msg("failed to parse environment variables")
	}

	settings = mergeSettings(settings, envVarSettings)
	if !settings.Valid() {
		logger.Error().Msg("missing required settings")
		return settings, fmt.Errorf("missing required settings")
	}
	return settings, nil
}
