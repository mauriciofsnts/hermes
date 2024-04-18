package config

import (
	"os"

	"github.com/ghodss/yaml"
)

var (
	Hermes *Config
)

func LoadConfigFromFile(configPath string) (*Config, error) {
	/* #nosec G304 */
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	return LoadConfigFromData(data)
}

func LoadConfigFromData(data []byte) (*Config, error) {
	var config Config

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	ensureNotNil(&config)

	config.AppsByAPIKey = make(map[string]*AppConfig)

	for _, app := range config.Apps {
		config.AppsByAPIKey[app.APIKey] = app
	}

	Hermes = &config
	return &config, nil
}

func ensureNotNil(cfg *Config) {
	if cfg.Log == nil {
		cfg.Log = &LogConfig{}
	}
	if cfg.Http == nil {
		cfg.Http = &HTTPConfig{}
	}
	if cfg.Apps == nil {
		cfg.Apps = make(map[string]*AppConfig)
	}
	// if cfg.Redict == nil {
	// 	cfg.Redict = &RedictConfig{}
	// }
	// if cfg.Public == nil {
	// 	cfg.Public = &AppConfig{}
	// }
}
