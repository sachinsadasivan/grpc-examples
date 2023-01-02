package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/sachinsadasivan/bidirectional-streaming-rpc-example/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultName  = "Bidirectional Streaming RPC Example"
	certFile     = "cert.pem"
	fileToStream = "SampleFile.txt"
)

var (
	addr = flag.String("addr", "localhost:50000", "the address to connect to")
	name = flag.String("name", defaultName, "Name of the application")
)

func main() {

	creds, err := credentials.NewClientTLSFromFile(certFile, "")

	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := messages.NewStreamingServiceClient(conn)

	stream, err := client.StreamData(context.Background())

	if err != nil {
		log.Fatalf("could not stream: %v", err)
	}

	doneCh := make(chan struct{})

	fileData := []byte{}
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				doneCh <- struct{}{}
				break
			}
			if err != nil {
				log.Fatalf("could not receive: %v", err)
			}

			fmt.Printf("Client: Recived data with lenght %d\n", len(fileData))

			mydir, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Current working directory: " + mydir)

			file := filepath.Join(mydir, "downloaded-"+fileToStream)
			err = os.WriteFile(file, fileData, 0666)
			if err != nil {
				log.Fatal(err)
			}

			fileData = append(fileData, in.Data...)

		}
	}()

	f, err := os.Open(fileToStream)
	fmt.Printf("Client: Streaming file: %s\n", fileToStream)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = stream.Send(&messages.FileStreamingRequest{})
	if err != nil {
		log.Fatalf("could not send: %v", err)
	}

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
		stream.Send(&messages.FileStreamingRequest{Data: chunk})
	}

	stream.CloseSend()
	<-doneCh

}
