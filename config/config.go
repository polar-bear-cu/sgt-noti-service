package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	SMTPHost     string
	SMTPPort     string
	SMTPFrom     string
	SMTPUsername string
	SMTPPassword string

	MongoURI    string
	MongoDBName string

	RabbitMQURL string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Port: env("PORT", "8080"),

		SMTPHost:     env("SMTP_HOST", "localhost"),
		SMTPPort:     env("SMTP_PORT", "1025"),
		SMTPFrom:     env("SMTP_FROM", "noreply@subglutee.local"),
		SMTPUsername: env("SMTP_USERNAME", ""),
		SMTPPassword: env("SMTP_PASSWORD", ""),

		MongoURI:    env("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: env("MONGO_DB_NAME", "noti"),

		RabbitMQURL: env("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
