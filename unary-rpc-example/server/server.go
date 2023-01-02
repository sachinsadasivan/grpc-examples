package main

import (
	"context"
	"log"
	"net"

	"github.com/sachinsadasivan/unary-rpc-example/messages"

	"google.golang.org/grpc"
)

const port = ":8085"

type server struct {
	messages.HelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *messages.HelloRequest) (*messages.HelloResponse, error) {
	log.Printf("Recived message from Client: %v", req.Name)
	return &messages.HelloResponse{Message: "Hello " + req.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	messages.RegisterHelloServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
