package main

import (
	"context"
	"fmt"
	pb "github.com/golang-grpc/greet/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func (s *Server) Greet(ctx context.Context, request *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("function was invoked %v\n", request.GetFirstName())
	return &pb.GreetResponse{
		Result: "Hello " + request.GetFirstName(),
	}, nil
}

func (s *Server) GreetMany(req *pb.GreetRequest, stream pb.GreetService_GreetManyServer) error {

	for i := 0; i < 10; i++ {

		res := fmt.Sprintf("Hello %s for %d", req.FirstName, i)

		err := stream.Send(&pb.GreetResponse{
			Result: res,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) LongGreet(stream pb.GreetService_LongGreetServer) error {

	res := ""

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.GreetResponse{
				Result: res,
			})
		}

		res += msg.FirstName + " "
	}

}

func (s *Server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error %v\n", err)
		}

		res := "Hello " + req.FirstName + " !"

		err = stream.Send(&pb.GreetResponse{
			Result: res,
		})

		if err != nil {
			log.Fatalf("error %v\n", err)
		}
	}

}

func (s *Server) GreetWithDeadline(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {

	for i := 0; i < 3; i++ {

		if ctx.Err() == context.DeadlineExceeded {
			log.Println("client has cancelled request")
			return nil, status.Error(codes.Canceled, "the client has cancelled the request")
		}

		time.Sleep(1 * time.Second)
	}

	return &pb.GreetResponse{Result: "Hello " + req.FirstName}, nil

}
