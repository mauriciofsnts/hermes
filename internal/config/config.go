package config

import "log/slog"

type Config struct {
	SMTP  SMTPConfig
	Http  *HTTPConfig
	Log   *LogConfig
	Redis *RedisConfig
	Kafka *KafkaConfig

	Apps         map[string]*AppConfig
	AppsByAPIKey map[string]*AppConfig `yaml:"-" json:"-"`
}

type RedisConfig struct {
	Address  string
	Password string
	Topic    string
}

type KafkaConfig struct {
	Address string
	Topic   string
	Brokers []string
}

type HTTPConfig struct {
	Port int
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
