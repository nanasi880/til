package main

import (
	"log"
	"net"

	"github.com/nanasi880/til/go/grpc/proto"
	"google.golang.org/grpc"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	helloService := new(HelloServiceServer)

	proto.RegisterHelloServiceServer(server, helloService)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
