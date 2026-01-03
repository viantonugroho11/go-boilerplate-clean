package config

import (
	"os"
	"strings"

)

type Configuration struct{
	App           App               `json:"app"`
	Database      PostgreDB         `json:"database"`
	KafkaProducer map[string]string `json:"kafka_producer"`
	Kafka         Kafka             `json:"kafka"`
	Consumers     Consumers         `json:"consumers"`
}

type AppConfig struct {
	App           App               `json:"app"`
	Database      PostgreDB         `json:"database"`
	KafkaProducer map[string]string `json:"kafka_producer"`
	Kafka         Kafka             `json:"kafka"`
	Consumers     Consumers         `json:"consumers"`
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

var Config *AppConfig = &AppConfig{}

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
