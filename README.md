# Devfest gRPC Examples

This project demonstrates simple gRPC implementations in Go, focusing on Unary and Bi-directional Streaming RPCs. Using Protocol Buffers for message serialization and gRPC as the communication framework, we’ll implement a basic service where a client sends a greeting request to a server, which responds with a personalized message.

## Prerequisites

### Installing gRPC Tooling

1. **Create a Go Project:**
   
```bash
mkdir devfest-grpc
cd devfest-grpc
go mod init devfest-grpc
```
2. **Install gRPC for Go:**
    
```bash
go get -u google.golang.org/grpc
```
3. **Install Protocol Buffer Tools for Go:**
```bash
go get -u google.golang.org/protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
### Defining Protocol Buffers
1. **Create a Folder for Message Definitions:**
   
```bash
mkdir messages
```
2. **Create a File for Protocol Buffers:** Create messages.proto inside the messages folder with the following content:
```bash
syntax = "proto3";
option go_package = "devfest-grpc/messages";

message HelloRequest {
    string Name = 1;
}

message HelloResponse {
    string Message = 1;
}

service HelloService {
    rpc SayHello (HelloRequest) returns (HelloResponse) {}
}

```
Explanation:

syntax specifies the Protocol Buffers version.
option go_package sets the Go package path.
HelloRequest and HelloResponse define the message types for requests and responses.
HelloService declares an RPC method SayHello, accepting a HelloRequest and returning a HelloResponse.

3. **Generate Code:** Use the protoc command to generate Go code from messages.proto:
   
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative messages/messages.proto
```
Explanation:

--go_out and --go-grpc_out specify output directories for standard and gRPC code.
paths=source_relative maintains the folder structure, simplifying imports in Go code.

### Server Implementation
1. **Create Server Folder and Go File:** Create server.go inside a server folder with the following content:
```bash
package main

import (
    "context"
    "log"
    "net"
    "devfest-grpc/messages"
    "google.golang.org/grpc"
)

const port = ":8085"

type server struct {
    messages.HelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *messages.HelloRequest) (*messages.HelloResponse, error) {
    log.Printf("Received message from Client: %v", req.Name)
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

```
Explanation:

SayHello: Implements the RPC method, logging the client's name and responding with a greeting.
main: Sets up the server to listen on port 8085 and registers the HelloServiceServer
### Server Implementation
1. **Create Client Folder and Go File:** Create client.go inside a client folder with the following content:
```bash
package main

import (
    "context"
    "flag"
    "log"
    "time"
    "devfest-grpc/messages"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

const (
    defaultName = "Unary RPC Example"
)

var (
    addr = flag.String("addr", "localhost:8085", "the address to connect to")
    name = flag.String("name", defaultName, "Name to greet")
)

func main() {
    flag.Parse()
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := messages.NewHelloServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.SayHello(ctx, &messages.HelloRequest{Name: *name})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("Greeting from Server: %s", r.GetMessage())
}

```
Explanation:

main: Connects to the gRPC server and sends a HelloRequest, logging the server's response.

### Running the Example
1. **Start the Server:**
```bash
go run server/server.go
```
**Expected output:**
```bash
2024/11/14 23:06:35 Server listening at [::]:8085
```
2. **Start the Client:**
```bash
go run client/client.go
```
**Expected output:**
```bash
2024/11/14 23:07:08 Greeting from Server: Hello Unary RPC Example
```
### Conclusion

This project demonstrates the basics of gRPC in Go through a simple Unary RPC service, covering setup, message definition, and client-server communication. The example serves as a starting point for building more complex services and exploring advanced gRPC features.












   
