package main

import (
	"context"
	"fmt"
	"helloGRPC/hello/hellopb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	fmt.Printf("Hello function line 15! %v \n", req)

	firstName := req.GetHello().GetFirstName()
	prefix := req.GetHello().GetPrefix()

	customHello := "Welcome ! " + prefix + " " + firstName

	res := &hellopb.HelloResponse{
		CustomHello: customHello,
	}
	return res, nil
}

func (*server) HelloManyLanguages(req *hellopb.HelloManyLanguagesRequest, stream hellopb.HelloService_HelloManyLanguagesServer) error {
	fmt.Printf("HelloManyLanguages invoked!! %v 🚀", req)
	langs := [3]string{"Salut! ", "Hola! ", "Hi! "}

	firstName := req.GetHello().GetFirstName()
	prefix := req.GetHello().GetPrefix()

	for _, helloLang := range langs {
		helloLanguage := helloLang + prefix + " " + firstName

		res := &hellopb.HelloManyLanguagesResponse{
			AnyHello: helloLanguage,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) HelloGoodBye(stream hellopb.HelloService_HelloGoodByeServer) error {
	fmt.Println("HelloGoodBye invoked!! 🚀")

	goodbye := "Adios perros! "

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			//Once the stream is finished we are going to send the response

			return stream.SendAndClose(&hellopb.HelloGoodByeResponse{Goodbye: goodbye})
		}

		if err != nil {
			log.Fatalf("Error reading the client stream %v", err)
		}

		firstName := req.GetHello().GetFirstName()
		prefix := req.GetHello().GetPrefix()

		goodbye += prefix + " " + firstName + " "

	}
}

func (*server) GoodByeBidirectional(stream hellopb.HelloService_GoodByeBidirectionalServer) error {
	fmt.Println("GoodByeBidirectional invoked!! 🚀")
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			//Once the stream is finished respond nil
			return nil
		}

		if err != nil {
			log.Fatalf("Error reading the client stream %v", err)
			return nil
		}

		firstName := req.GetHello().GetFirstName()
		prefix := req.GetHello().GetPrefix()

		goodbye := "Goodbye " + prefix + " " + firstName + " "

		sendErr := stream.Send(&hellopb.GoodByeBidirectionalResponse{
			Goodbye: goodbye,
		})

		if sendErr != nil {
			log.Fatalf("Error reading the client stream %v", err)
		}
	}
}

func main() {
	fmt.Println("Go server running!!🚀")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	hellopb.RegisterHelloServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
