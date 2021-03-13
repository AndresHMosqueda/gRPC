package main

import (
	"context"
	"fmt"
	"helloGRPC/hello/hellopb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Go client is running 🚀")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}

	defer cc.Close()
	c := hellopb.NewHelloServiceClient(cc)

	// helloUnary(c)
	// helloServerStreaming(c)
	// goodbyeClientStreaming(c)
	goodbyeBidirectionalStreaming(c)
}

func helloUnary(c hellopb.HelloServiceClient) {
	fmt.Println("Starting Unary RPC")

	req := &hellopb.HelloRequest{
		Hello: &hellopb.Hello{
			FirstName: "Andres",
			Prefix:    "Mr",
		},
	}

	res, err := c.Hello(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling grpc: \n %v", err)
	}

	fmt.Printf("Response: %v", res.CustomHello)
}

func helloServerStreaming(c hellopb.HelloServiceClient) {
	fmt.Println("Starting Server Streaming RPC")

	req := &hellopb.HelloManyLanguagesRequest{
		Hello: &hellopb.Hello{
			FirstName: "Andre",
			Prefix:    "Mister",
		},
	}

	restStream, err := c.HelloManyLanguages(context.Background(), req)

	if err != nil {
		log.Fatalf("Error calling HelloManyLanguages: \n %v", err)
	}

	for {
		msg, err := restStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}

		fmt.Printf("Response: %v \n", msg.GetAnyHello())
	}
}

func goodbyeClientStreaming(c hellopb.HelloServiceClient) {
	fmt.Println("Starting Client Streaming RPC")

	requests := []*hellopb.HelloGoodByeRequest{
		{
			Hello: &hellopb.Hello{
				FirstName: "Andre Marin!",
				Prefix:    "Mister",
			},
		},
		{
			Hello: &hellopb.Hello{
				FirstName: "Cura!",
				Prefix:    "Sr",
			},
		},
		{
			Hello: &hellopb.Hello{
				FirstName: "Maria!",
				Prefix:    "Srita",
			},
		},
		{
			Hello: &hellopb.Hello{
				FirstName: "Marcelino!",
				Prefix:    "Mister",
			},
		},
	}

	stream, err := c.HelloGoodBye(context.Background())
	if err != nil {
		log.Fatalf("Error reading stream: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending requests... %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Microsecond)

	}

	goodbye, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error reading stream: %v", err)
	}

	fmt.Println("Respuesta: ", goodbye)

}

func goodbyeBidirectionalStreaming(c hellopb.HelloServiceClient) {
	fmt.Println("Starting BidirectionalStreaming RPC")

	//create stream to call the server
	stream, err := c.GoodByeBidirectional(context.Background())

	requests := []*hellopb.GoodByeBidirectionalRequest{
		{
			Hello: &hellopb.Hello{
				FirstName: "Andre Marin!",
				Prefix:    "Mister",
			},
		},
		{
			Hello: &hellopb.Hello{
				FirstName: "Cura!",
				Prefix:    "Sr",
			},
		},
		{
			Hello: &hellopb.Hello{
				FirstName: "Maria!",
				Prefix:    "Srita",
			},
		},
		{
			Hello: &hellopb.Hello{
				FirstName: "Marcelino!",
				Prefix:    "Mister",
			},
		},
	}
	if err != nil {
		log.Fatalf("Error creating the stream: %v", err)
	}

	waitChannel := make(chan struct{})
	//send many messages to the server (go routines)
	go func() {
		for _, req := range requests {
			log.Printf("Sending messages... %v", req)
			stream.Send(req)
			time.Sleep(300 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//receive messages from the server (go routines)
	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error receiving stream %v", err)
				break
			}

			fmt.Printf("Respuesta: %v\n", res.GetGoodbye())
		}
		close(waitChannel)
	}()

	//block when everything is completed or closed
	<-waitChannel
}
