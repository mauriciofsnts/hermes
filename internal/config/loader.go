package config

import (
	"os"
	"strconv"

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

	// Override with environment variables
	applyEnvOverrides(&config)

	config.AppsByAPIKey = make(map[string]*AppConfig)

	for _, app := range config.Apps {
		config.AppsByAPIKey[app.APIKey] = app
	}

	// Validar configuração
	if err := ValidateConfig(&config); err != nil {
		return nil, err
	}

	Hermes = &config
	return &config, nil
}

// applyEnvOverrides allows environment variables to override config.yaml values
func applyEnvOverrides(cfg *Config) {
	// SMTP Configuration
	if host := os.Getenv("SMTP_HOST"); host != "" {
		cfg.SMTP.Host = host
	}
	if port := os.Getenv("SMTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.SMTP.Port = p
		}
	}
	if username := os.Getenv("SMTP_USERNAME"); username != "" {
		cfg.SMTP.Username = username
	}
	if password := os.Getenv("SMTP_PASSWORD"); password != "" {
		cfg.SMTP.Password = password
	}
	if from := os.Getenv("SMTP_FROM"); from != "" {
		cfg.SMTP.Sender = from
	}

	// Redis Configuration
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		cfg.Redis.Address = addr
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}
	if topic := os.Getenv("REDIS_TOPIC"); topic != "" {
		cfg.Redis.Topic = topic
	}

	// HTTP Configuration
	if port := os.Getenv("HTTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Http.Port = p
		}
	}

	// Database Configuration
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.PG.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if p, err := strconv.Atoi(dbPort); err == nil {
			cfg.PG.Port = p
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.PG.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.PG.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.PG.DBName = dbName
	}
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
	if cfg.PG == nil {
		cfg.PG = &PGConfig{}
	}
}
