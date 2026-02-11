package config

import (
	"crypto-aggregator-service/internal/models"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
)

// Configurations Application wide configurations
type Configurations struct {
	Server ServerConfigurations `koanf:"server"`
	App    AppConfigurations    `koanf:"app"`
	Keys   KeysConfigurations   `koanf:"keys"`
}

// ServerConfigurations Server configurations
type ServerConfigurations struct {
	Port            int `koanf:"port"`
	RefreshInterval int `koanf:"refresh_interval"`
}

// AppConfigurations App configurations
type AppConfigurations struct {
	Layout []ItemConfig `koanf:"layout"`
}

// KeysConfigurations asymmetric keys
type KeysConfigurations struct {
	Public string `koanf:"public"`
}

// ItemConfig represents a row in config.json.
// It maps to the domain component but adds the necessary "Vendor" config.
type ItemConfig struct {
	ID        int    `json:"id"`
	Component string `json:"component"`
	Vendor    string `json:"vendor"` // Configuration only!
}

// LoadConfig Loads configurations depending upon the environment
func LoadConfig(logger *zap.SugaredLogger) *Configurations {
	k := koanf.New(".")
	err := k.Load(file.Provider("resources/config.yaml"), yaml.Parser())
	if err != nil {
		logger.Fatalf("Failed to locate configurations. %v", err)
	}

	// Searches for env variables and will transform them into koanf format
	// e.g. SERVER_PORT variable will be server.port: value
	err = k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil)
	if err != nil {
		logger.Fatalf("Failed to replace environment variables. %v", err)
	}

	var configuration Configurations

	err = k.Unmarshal("", &configuration)
	if err != nil {
		logger.Fatalf("Failed to load configurations. %v", err)
	}

	return &configuration
}

// Helper to convert Config -> Domain
func (c *AppConfigurations) ToDomain() []models.Component {
	list := make([]models.Component, len(c.Layout))
	for i, item := range c.Layout {
		list[i] = models.Component{
			ID:        item.ID,
			Component: models.ComponentType(item.Component),
			Model:     nil, // Starts empty
		}
	}
	return list
}

// Helper to extract Vendor Map (ID -> Vendor)
func (c *AppConfigurations) GetVendorMap() map[int]string {
	m := make(map[int]string)
	for _, item := range c.Layout {
		m[item.ID] = item.Vendor
	}
	return m
}
