package main

import (
	"context"
	"flag"
	"io"
	"time"

	pb "gesture-project-aws-grpc-client/message"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "myClientName"
)

var (
	addr = flag.String("addr", "192.168.0.6:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	initLogger()

	log.Info("Started...")

	flag.Parse()

	ch := createChannel()
	defer ch.Close()

	for {
		log.Infof("Attempting to connect to [%s] with name [%s]", *addr, *name)

		// Set up a connection to the server.
		conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("did not connect: %v", err)
		} else {
			defer conn.Close()
			c := pb.NewGreeterClient(conn)

			// Contact the server and print out its response.
			// ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
			if err != nil {
				log.Errorf("could not greet: %v", err)
			} else {
				log.Infof("Connected to [%s] with name [%s]", *addr, *name)

				counter := 0
				for {
					data, err := stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Infof("Failed to receive a message from server: %v", err)
						counter++
						if counter > 50 {
							break
						}
					}
					if data == nil {
						continue
					}

					counter = 0
					// log.Infof("Greeting: %s, %s", data.GetKey(), data.GetValue())

					sendMessage(ch, data.GetKey(), data.GetValue())
				}
			}
		}

		log.Notice("Connection closed by server. Reconnecting in 1s...")
		// sleep 1 second
		time.Sleep(time.Second * 1)
	}
}
