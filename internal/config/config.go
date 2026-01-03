package config

import (
	"fmt"
	"strings"

)

type Configuration struct {
	App           App               `json:"app"`
	Database      PostgreDB         `json:"database"`
	KafkaProducer map[string]string `json:"kafka_producer"`
	Kafka         Kafka             `json:"kafka"`
	Redis         Redis             `json:"redis"`
	// Consumers     Consumers         `json:"consumers"`
}


func (c Configuration) PGDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.DBName, c.Database.SSLMode)
}


func (c Configuration) KafkaBrokersList() []string {
	parts := strings.Join(c.Kafka.Brokers, ",")
	return strings.Split(parts, ",")
}
