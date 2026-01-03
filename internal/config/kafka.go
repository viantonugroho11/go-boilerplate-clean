package config

type Kafka struct {
	Brokers []string `json:"brokers"`
	ClientID string `json:"client_id"`
	GroupID string `json:"group_id"`
	Topic string `json:"topic"`
}