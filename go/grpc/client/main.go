package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nanasi880/til/go/grpc/proto"
	"google.golang.org/grpc"
)

// https://christina04.hatenablog.com/entry/2017/11/13/190000
func main() {

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	req := proto.HelloRequest{
		Name: "sample client",
	}

	resp, err := client.Hello(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Message)
}
