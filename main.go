package main

import (
	"time"
)

func main() {
	initLogger()

	log.Info("Started...")

	ch := createChannel()
	defer ch.Close()

	consumer := createConsumer()
	defer consumer.Close()

	readMessages(consumer, func(logEntry LogEntry) {
		log.Info("Processing message...")
		log.Info(logEntry)

		pinCode := logEntry.Key

		// get timestamp local time with timezone
		timestamp := time.Now().Format("2006-01-02 15:04:05.999-07:00")
		message := map[string]interface{}{
			"timestamp": timestamp,
			"value":     logEntry.Value,
		}

		sendMessage(ch, pinCode, message)
	})

}
