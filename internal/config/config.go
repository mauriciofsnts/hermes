package config

type Config struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
	DefaultFrom  string
}
