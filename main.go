package main

// Task: Write HTTP request code to pull latest message from Kafka topic "Gestures"
// Reference: Consume Messages - Kafka API: https://docs.upstash.com/kafka/kafkaapi
import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"

	"encoding/json"
)

func printSomething(args string, val int) {
	fmt.Println(args, val)
}
// layout of the printed message:
type LogEntry struct {
	Timestamp string `json:"Time"`
	Value     map[string]string `json:"Value"`
}

func main() {
	fmt.Println("Hello, World!")
	printSomething("Hello", 1)

	// New main from documentation:
	mechanism, err := scram.Mechanism(scram.SHA256, 
		"d2lubmluZy1jaGlwbXVuay0xMzI4NiRu6Cm3A9Mo_Q6mThRD_7s0zqgOo3T7pIE", "9f10c92e53ff4e96baafdadbc2c9c6fe")
	if err != nil {
	//   log.Fatalln(err)
	  log.Println("Error: ", err)
	}
	
	dialer := &kafka.Dialer{
	  SASLMechanism: mechanism,
	  TLS:           &tls.Config{},
	}
	
	r := kafka.NewReader(kafka.ReaderConfig{
	  Brokers:  []string{"winning-chipmunk-13286-us1-kafka.upstash.io:9092"},
	  GroupID: "$GROUP_NAME",
	  Topic:   "Gestures",
	  Dialer:  dialer,
	})
	defer r.Close()
	
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
	
		m, err := r.ReadMessage(ctx)
		if err != nil {
		//   log.Fatalln(err)
		log.Println("Error: ", err)
		}
		
		var logEntry LogEntry
		
		err_ := json.Unmarshal(m.Value, &logEntry)
		if err_ != nil {
			fmt.Println("Error parsing log entry:", err_)
			return
		}
		
		x := logEntry.Value["x"]
		y := logEntry.Value["y"]
		
		// fmt.Println("Type of m:", reflect.TypeOf(m))

		log.Println("the raw output:")
		log.Printf("%+v\n", m)
		log.Println("the parsed output: ")
		fmt.Println("x:", x)
		fmt.Println("y:", y)
	}
}
