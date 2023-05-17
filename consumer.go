package main

// Task: Write HTTP request code to pull latest message from Kafka topic "Gestures"
// Reference: Consume Messages - Kafka API: https://docs.upstash.com/kafka/kafkaapi
import (
	"context"
	"crypto/tls"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

// layout of the printed message:
type LogEntry struct {
	Key     string `json:"key"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Gesture string `json:"gesture"`
}

func createConsumer() *kafka.Reader {
	// New main from documentation:
	mechanism, err := scram.Mechanism(scram.SHA256,
		"d2lubmluZy1jaGlwbXVuay0xMzI4NiRu6Cm3A9Mo_Q6mThRD_7s0zqgOo3T7pIE", "9f10c92e53ff4e96baafdadbc2c9c6fe")
	if err != nil {
		log.Fatalf("Error: ", err)
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
		TLS:           &tls.Config{},
	}

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"winning-chipmunk-13286-us1-kafka.upstash.io:9092"},
		//   GroupID: "myGroupId",
		Topic:       "Gestures",
		StartOffset: kafka.LastOffset,
		Dialer:      dialer,
	})
	// defer consumer.Close()

	return consumer
}

func readMessages(consumer *kafka.Reader, processMsg func(LogEntry)) {
	for {
		log.Info("Waiting for message...")

		// ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		// defer cancel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		msg, err := consumer.ReadMessage(ctx)
		if err != nil {
			log.Error("Error:", err)
			continue
		}

		log.Info("Hi")

		// check if msg.Time is within the last 5 seconds
		if time.Since(msg.Time) > 5*time.Second {
			log.Warning("Message is too old, skipping...")
			continue
		}

		var logEntry LogEntry
		err_ := json.Unmarshal(msg.Value, &logEntry)
		if err_ != nil {
			log.Error("Error parsing log entry:", err_)
			continue
		}
		logEntry.Key = string(msg.Key)

		log.Info("the raw output:")
		log.Infof("%+v\n", string(msg.Value))
		// log.Info("rcvd:", string(msg.Value))
		// log.Info("logEntry:", logEntry)

		processMsg(logEntry)
	}
}
