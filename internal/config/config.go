package config

type Config struct {
	SmtpHost      string
	SmtpPort      int
	SmtpUsername  string
	SmtpPassword  string
	DefaultFrom   string
	AllowedOrigin string
	Redis         Redis
	Kafka         Kafka
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
