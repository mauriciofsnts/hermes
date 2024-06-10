package config

import "log/slog"

type Config struct {
	SMTP  SMTPConfig
	PG    *PGConfig
	Http  *HTTPConfig
	Log   *LogConfig
	Redis *RedisConfig

	Apps         map[string]*AppConfig
	AppsByAPIKey map[string]*AppConfig `yaml:"-" json:"-"`
}

type RedisConfig struct {
	Address  string
	Password string
	Topic    string
}

type HTTPConfig struct {
	Port           int
	AllowedOrigins []string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type AppConfig struct {
	Enabled           bool
	APIKey            string
	AllowedOrigins    []string
	LimitPerIPPerHour int
	Discord           *DiscordWebhook
}

type DiscordWebhook struct {
	Token string
	ID    string
}

type LogType string

const (
	LogTypeText    LogType = "text"
	LogTypeJSON    LogType = "json"
	LogTypeColored LogType = "colored"
)

type LogConfig struct {
	Level      slog.Level
	Type       LogType
	ShowSource bool
}

type PGConfig struct {
	Migrate  bool
	Host     string
	Port     int
	User     string
	Password string
	SSLMode  string
	DBName   string
}
