package config

type Config struct {
	DefaultFrom    string
	AllowedOrigins []string
	Redis          Redis
	Kafka          Kafka
	Smtp           Smtp
	RateLimit      int
}

type Redis struct {
	Password string
	Host     string
	Port     int
	Topic    string
	Enabled  bool
}

type Kafka struct {
	Host    string
	Port    int
	Enabled bool
	Topic   string
	Brokers []string
}

type Smtp struct {
	Host     string
	Port     int
	Username string
	Password string
}
