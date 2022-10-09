package main

import (
	pb "github.com/golang-grpc/calculator/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

var address = "0.0.0.0:8081"

type Server struct {
	pb.CalulatorServiceServer
}

func main() {

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	s := grpc.NewServer()

	pb.RegisterCalulatorServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("error %v\n", err)
	}
}
