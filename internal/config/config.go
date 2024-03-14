package config

type Config struct {
	DefaultFrom string
	Hermes      Hermes
	Redis       Redis
	Kafka       Kafka
	SMTP        SMTP
}

type Hermes struct {
	RateLimit int
	Apikeys   []string
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

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
}
