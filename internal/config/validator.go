package config

import (
	"errors"
	"fmt"
)

// ValidateConfig verifica se todos os campos obrigatórios estão configurados
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	// Validar SMTP
	if err := validateSMTP(&cfg.SMTP); err != nil {
		return fmt.Errorf("invalid SMTP config: %w", err)
	}

	// Validar HTTP
	if err := validateHTTP(cfg.Http); err != nil {
		return fmt.Errorf("invalid HTTP config: %w", err)
	}

	// Validar Apps
	if len(cfg.Apps) == 0 {
		return errors.New("at least one app must be configured")
	}

	for name, app := range cfg.Apps {
		if app.APIKey == "" {
			return fmt.Errorf("app '%s' must have an API key", name)
		}
		if len(app.EnabledFeatures) == 0 {
			return fmt.Errorf("app '%s' must have at least one enabled feature", name)
		}
	}

	// Validar Redis se estiver configurado
	if cfg.Redis != nil && cfg.Redis.Address != "" {
		if err := validateRedis(cfg.Redis); err != nil {
			return fmt.Errorf("invalid Redis config: %w", err)
		}
	}

	return nil
}

func validateSMTP(smtp *SMTPConfig) error {
	if smtp == nil {
		return errors.New("SMTP config is required")
	}

	if smtp.Host == "" {
		return errors.New("SMTP host is required")
	}

	if smtp.Port == 0 {
		return errors.New("SMTP port is required")
	}

	if smtp.Port < 1 || smtp.Port > 65535 {
		return fmt.Errorf("SMTP port must be between 1 and 65535, got %d", smtp.Port)
	}

	if smtp.Username == "" {
		return errors.New("SMTP username is required")
	}

	if smtp.Password == "" {
		return errors.New("SMTP password is required")
	}

	return nil
}

func validateHTTP(http *HTTPConfig) error {
	if http == nil {
		return errors.New("HTTP config is required")
	}

	if http.Port == 0 {
		return errors.New("HTTP port is required")
	}

	if http.Port < 1 || http.Port > 65535 {
		return fmt.Errorf("HTTP port must be between 1 and 65535, got %d", http.Port)
	}

	return nil
}

func validateRedis(redis *RedisConfig) error {
	if redis.Address == "" {
		return errors.New("Redis address is required when configured")
	}

	if redis.Topic == "" {
		return errors.New("Redis topic is required when configured")
	}

	return nil
}
