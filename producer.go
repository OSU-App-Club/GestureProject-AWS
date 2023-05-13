package main

import (
	"encoding/json"
	"strings"

	"github.com/streadway/amqp"
)

var exchangeName = "msgExchange"

func createChannel() *amqp.Channel {
	// Create a connection factory with the URI of your Amazon MQ RabbitMQ broker
	conn, err := amqp.Dial("amqps://pdohatix:7I0Wco0RKRe8z4LC3uGx4k-7uqWqNCmL@beaver.rmq.cloudamqp.com/pdohatix")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	// defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	// defer ch.Close()

	// Declare the exchange
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
			routingKey := strings.Split(ret.RoutingKey, "-")[1]
			message := string(ret.Body)
			statusCode := ret.ReplyCode
			reason := ret.ReplyText

			log.Warning("[RabbitMQ] - Message returned:", routingKey, "-", message)
			log.Warning("[RabbitMQ] - Status code:", statusCode)
			log.Warning("[RabbitMQ] - Reason:", reason)
		}
	}()

	return ch
}

func sendMessage(ch *amqp.Channel, pinCode string, message map[string]interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		log.Error("Failed to marshal message:", err)
		return err
	}

	// Create a message properties object with the desired headers and properties
	properties := amqp.Publishing{
		ContentType:  "text/plain",
		Body:         body,
		DeliveryMode: 2,
	}

	// Create a routing key
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
		log.Error("Failed to publish a message:", err)
		return err
	}

	log.Info("Sent message:", string(body))
	return nil
}
