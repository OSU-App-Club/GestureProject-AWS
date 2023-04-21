package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	// Create a connection factory with the URI of your Amazon MQ RabbitMQ broker
	conn, err := amqp.Dial("amqps://pdohatix:7I0Wco0RKRe8z4LC3uGx4k-7uqWqNCmL@beaver.rmq.cloudamqp.com/pdohatix")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare the exchange
	exchangeName := "msgExchange"
	err = ch.ExchangeDeclare(
		exchangeName, // exchange name
		"direct",     // exchange type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	// Add a callback for when the message is returned (e.g. if there is no queue bound to the exchange with the specified routing key)
	retCh := ch.NotifyReturn(make(chan amqp.Return, 1))
	go func() {
		for ret := range retCh {
			message := string(ret.Body)
			statusCode := ret.ReplyCode
			reason := ret.ReplyText

			log.Printf("Message returned: %s\n", message)
			log.Printf("Status code: %d\n", statusCode)
			log.Printf("Reason: %s\n", reason)
		}
	}()

	// create infinite loop
	for {
		// Create a message and convert it to a byte array
		// create random value
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		value := rnd.Intn(100)
		// get timestamp local time with timezone
		timestamp := time.Now().Format("2006-01-02 15:04:05.999-07:00")
		message := map[string]interface{}{
			"Name":      "John",
			"Age":       value,
			"Timestamp": timestamp,
		}
		body, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Failed to marshal message: %s", err)
		}

		// Create a message properties object with the desired headers and properties
		properties := amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			DeliveryMode: 2,
		}

		// Create a routing key
		pinCode := "123456"
		routingKey := "key-" + pinCode

		// Publish the message to the exchange with the desired routing key
		const mandatory = true
		err = ch.Publish(
			exchangeName, // exchange name
			routingKey,   // routing key
			mandatory,    // mandatory
			false,        // immediate
			properties,   // message properties
		)
		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}

		log.Printf("Sent message: %v\n", message)

		// Sleep for 5 seconds
		time.Sleep(5 * time.Second)
	}
}