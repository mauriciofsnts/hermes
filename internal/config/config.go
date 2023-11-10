package config

type Config struct {
	DefaultFrom   string
	AllowedOrigin string
	Redis         Redis
	Kafka         Kafka
	Smtp          Smtp
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
