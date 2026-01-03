package config

import (
	"os"
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
}

func Load() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DATABASE_URL")
	return AppConfig{
		Port:        port,
		DatabaseURL: dbURL,
		DBHost:      getEnvDefault("DB_HOST", "127.0.0.1"),
		DBPort:      getEnvDefault("DB_PORT", "5432"),
		DBUser:      getEnvDefault("DB_USER", "postgres"),
		DBPassword:  getEnvDefault("DB_PASSWORD", ""),
		DBName:      getEnvDefault("DB_NAME", "appdb"),
		DBSSLMode:   getEnvDefault("DB_SSLMODE", "disable"),
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
