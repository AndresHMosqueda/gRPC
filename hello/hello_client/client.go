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
	goodbyeClientStreaming(c)
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
