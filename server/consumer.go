package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	bootstrapServers = "pkc-rgm37.us-west-2.aws.confluent.cloud:9092"
	APIKey           = "44TSILYQLLH3SNPX"
	APISecret        = "9g7KBW0emCSUXjb5Fa069DPqGrrTHhOf5UY78oZlZmLns1lLbAvnlQdBnu6rle/Y"
)

// layout of the printed message:
type LogEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func createConsumer() *kafka.Consumer {
	// Now consumes the record and print its value...
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     APIKey,
		"sasl.password":     APISecret,
		"group.id":          "my-groupd",
		"auto.offset.reset": "latest"})

	if err != nil {
		panic(fmt.Sprintf("Failed to create consumer: %s", err))
	}

	topic := "Gestures"
	topics := []string{topic}
	consumer.SubscribeTopics(topics, nil)
	// defer consumer.Close()

	return consumer
}

func readMessages(consumer *kafka.Consumer, processMsg func(LogEntry)) {
	for {
		// log.Info("Waiting for message...")

		message, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Error("Error:", err)
			continue
		}

		// check if message.Time is within the last 5 seconds
		if time.Since(message.Timestamp) > 5*time.Second {
			// log.Warning("Message is too old, skipping...")
			continue
		}

		var logEntry LogEntry
		logEntry.Key = string(message.Key)
		logEntry.Value = string(message.Value)

		// log.Infof("Message received: %s", string(message.Value))

		processMsg(logEntry)
	}
}
