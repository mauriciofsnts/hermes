package config

import (
	"testing"
)

func TestValidateConfigSuccess(t *testing.T) {
	cfg := &Config{
		SMTP: SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "user@example.com",
			Password: "password",
		},
		Http: &HTTPConfig{
			Port: 8080,
		},
		Apps: map[string]*AppConfig{
			"app1": {
				APIKey:          "key123",
				EnabledFeatures: []Feats{"email"},
			},
		},
	}

	err := ValidateConfig(cfg)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestValidateConfigNil(t *testing.T) {
	err := ValidateConfig(nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}
}

func TestValidateConfigMissingSMTPHost(t *testing.T) {
	cfg := &Config{
		SMTP: SMTPConfig{
			Port:     587,
			Username: "user",
			Password: "pass",
		},
		Http: &HTTPConfig{Port: 8080},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Error("Expected error for missing SMTP host")
	}
}

func TestValidateConfigInvalidPort(t *testing.T) {
	cfg := &Config{
		SMTP: SMTPConfig{
			Host:     "smtp.example.com",
			Port:     70000,
			Username: "user",
			Password: "pass",
		},
		Http: &HTTPConfig{Port: 8080},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Error("Expected error for invalid SMTP port")
	}
}

func TestValidateConfigNoApps(t *testing.T) {
	cfg := &Config{
		SMTP: SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "user",
			Password: "pass",
		},
		Http: &HTTPConfig{Port: 8080},
		Apps: map[string]*AppConfig{},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Error("Expected error for no apps configured")
	}
}

func TestValidateSMTPMissingPassword(t *testing.T) {
	smtp := &SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user",
		Password: "",
	}

	err := validateSMTP(smtp)
	if err == nil {
		t.Error("Expected error for missing SMTP password")
	}
}

func TestValidateHTTPMissingPort(t *testing.T) {
	http := &HTTPConfig{
		Port: 0,
	}

	err := validateHTTP(http)
	if err == nil {
		t.Error("Expected error for missing HTTP port")
	}
}
