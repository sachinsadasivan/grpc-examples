package main

import (
	"fmt"
	"io"

	"log"
	"net"
	"os"

	"github.com/sachinsadasivan/bidirectional-streaming-rpc-example/messages"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port         = ":50000"
	fileToStream = "job.log"
	certFile     = "cert.pem"
	keyFile      = "key.pem"
)

type server struct {
	messages.StreamingServiceServer
}

func (c *server) StreamData(stream messages.StreamingService_StreamDataServer) error {
	fmt.Println("Streaming data")

	doneCh := make(chan struct{})
	f, err := os.Open(fileToStream)
	fmt.Printf("Server: Streaming file %s\n", fileToStream)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	go func() {
		for {
			chunk := make([]byte, 64*1024)
			n, err := f.Read(chunk)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if n < len(chunk) {
				chunk = chunk[:n]
			}
			stream.Send(&messages.FileStreamingResponse{Data: chunk})
		}

		<-doneCh

	}()

	fileData := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("Server: File recived with lenght: %d\n", len(fileData))
			return stream.Send(&messages.FileStreamingResponse{Data: fileData})
		}
		if err != nil {
			return err
		}

		fmt.Printf("Server: Recived data with lenght %d\n", len(data.Data))

		fileData = append(fileData, data.Data...)

	}

}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)

	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}

	s := grpc.NewServer(opts...)

	messages.RegisterStreamingServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
