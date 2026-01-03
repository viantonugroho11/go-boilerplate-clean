package config

import (
	"os"
	"strings"
)

type AppConfig struct {
	Port        string
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	// Kafka
	KafkaBrokers  string
	KafkaClientID string
	KafkaGroupID  string
	KafkaTopic    string
	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       string
}

func Load() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DATABASE_URL")
	return AppConfig{
		Port:          port,
		DatabaseURL:   dbURL,
		DBHost:        getEnvDefault("DB_HOST", "127.0.0.1"),
		DBPort:        getEnvDefault("DB_PORT", "5432"),
		DBUser:        getEnvDefault("DB_USER", "postgres"),
		DBPassword:    getEnvDefault("DB_PASSWORD", ""),
		DBName:        getEnvDefault("DB_NAME", "appdb"),
		DBSSLMode:     getEnvDefault("DB_SSLMODE", "disable"),
		KafkaBrokers:  getEnvDefault("KAFKA_BROKERS", "127.0.0.1:9092"),
		KafkaClientID: getEnvDefault("KAFKA_CLIENT_ID", "go-boilerplate-clean"),
		KafkaGroupID:  getEnvDefault("KAFKA_GROUP_ID", "go-boilerplate-clean-group"),
		KafkaTopic:    getEnvDefault("KAFKA_TOPIC", "user-events"),
		RedisAddr:     getEnvDefault("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       getEnvDefault("REDIS_DB", "0"),
	}
}

func (c AppConfig) PGDSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	// pgx style DSN
	// example: postgres://user:pass@host:port/dbname?sslmode=disable
	user := c.DBUser
	pass := c.DBPassword
	host := c.DBHost
	port := c.DBPort
	db := c.DBName
	ssl := c.DBSSLMode
	if pass != "" {
		return "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + db + "?sslmode=" + ssl
	}
	return "postgres://" + user + "@" + host + ":" + port + "/" + db + "?sslmode=" + ssl
}

func getEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (c AppConfig) KafkaBrokersList() []string {
	parts := strings.Split(c.KafkaBrokers, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
