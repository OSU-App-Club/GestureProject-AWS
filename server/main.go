package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"

	pb "gesture-project-aws-grpc-server/message"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")

	// 1. Create a channel for the input, with a buffer size of 100.
	input = make(chan *pb.HelloReply, 100)
	// 2. Create a multiplexer, with the input channel as the argument.
	//   The multiplexer will listen to the input channel and distribute the messages to the subscribers.
	mux = NewMultiplexer(input)
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(in *pb.HelloRequest, stream pb.Greeter_SayHelloServer) error {
	log.Infof("New client connected: %v", in.GetName())

	output := make(chan *pb.HelloReply, 100)
	mux.Subscribe(output)

	counter := 0

	for message := range output {
		if err := stream.Send(message); err != nil {
			// log.Error("Error sending to gRPC client:", err)
			counter++

			if counter == 50 {
				log.Errorf("Disconnecting client [%s] b/c of 50 consecutive errors. Exiting with error: %s.", in.GetName(), err)
				mux.Unsubscribe(output)
				return err
			}
		} else {
			counter = 0
		}
	}

	log.Noticef("Client [%s] disconnected", in.GetName())
	mux.Unsubscribe(output)
	return nil
}

func main() {
	initLogger()

	log.Info("Started...")

	listenForKafkaMessages()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Infof("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func listenForKafkaMessages() {
	// run in background thread (go routine)
	go func() {
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
			// convert to string
			messageString, err := json.Marshal(message)
			if err != nil {
				log.Error("Failed to convert to string:", err)
				// skip this message
				return
			}

			data := &pb.HelloReply{
				Key:   pinCode,
				Value: string(messageString),
			}

			// handle the message
			input <- data
			log.Info("Message sent to clients.")

		})
	}()
}
