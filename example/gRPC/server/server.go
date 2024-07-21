package main

import (
	"context"
	"github.com/ahang7/go-IAM/example/gRPC/hellogRPC"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":50051"

type server struct {
	hellogRPC.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *hellogRPC.HelloRequest) (*hellogRPC.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &hellogRPC.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return
	}
	s := grpc.NewServer()
	hellogRPC.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
