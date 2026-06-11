package events

import (
	"log"
	"os"
	"strings"
)

// KafkaConfigured reports whether order-service should publish to Kafka.
// Set KAFKA_ENABLED=false to force-disable even when KAFKA_BROKERS is set.
func KafkaConfigured() bool {
	if strings.EqualFold(os.Getenv("KAFKA_ENABLED"), "false") {
		return false
	}
	return strings.TrimSpace(os.Getenv("KAFKA_BROKERS")) != ""
}

// LogKafkaMode prints whether the producer or noop publisher will be used.
func LogKafkaMode(enabled bool) {
	if enabled {
		brokers := os.Getenv("KAFKA_BROKERS")
		topic := os.Getenv("KAFKA_TOPIC")
		if topic == "" {
			topic = "order.placed"
		}
		log.Printf("kafka: enabled broker=%s topic=%s", brokers, topic)
		return
	}
	log.Println("kafka: disabled (set KAFKA_BROKERS to publish order.placed events)")
}
