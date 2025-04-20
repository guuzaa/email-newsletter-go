package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	DatabaseSettings `yaml:"database"`
	ApplicationPort  uint16 `yaml:"application_port" default:"8080"`
}

type DatabaseSettings struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Port         uint16 `yaml:"port"`
	Host         string `yaml:"host"`
	DatabaseName string `yaml:"database_name"`
}

func (setting Settings) PostgresSQLDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", setting.Host, setting.Username, setting.Password, setting.DatabaseName, setting.Port)
}

func (setting Settings) Address() string {
	return fmt.Sprintf("%s:%d", setting.Host, setting.ApplicationPort)
}

func Configuration(path string) (Settings, error) {
	var settings Settings
	// Load configuration from file or environment variables
	data, err := os.ReadFile(path)
	if err != nil {
		return settings, err
	}
	if err = yaml.Unmarshal(data, &settings); err != nil {
		return settings, err
	}
	return settings, nil
}
