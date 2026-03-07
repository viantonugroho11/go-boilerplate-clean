// Package broker berisi implementasi infrastruktur message broker.
//
// Kafka consumer dan producer memakai github.com/viantonugroho11/go-lib/kafka:
// - Consumer: dijalankan di transport/event/kafka (EventHandler + usecase).
// - Producer: wrapper typed producer di broker/kafka (UserEventPublisherKafka).
package broker
