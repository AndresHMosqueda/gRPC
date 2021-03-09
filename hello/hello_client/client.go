package main

import (
	"context"
	"fmt"
	"helloGRPC/hello/hellopb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Go client is running")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}

	defer cc.Close()
	c := hellopb.NewHelloServiceClient(cc)

	helloUnary(c)
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
