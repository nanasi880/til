package main

import (
	"context"
	"fmt"

	"github.com/nanasi880/til/go/grpc/proto"
)

type HelloServiceServer struct {
}

func (*HelloServiceServer) Hello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}
