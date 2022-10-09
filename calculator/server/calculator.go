package main

import (
	"context"
	pb "github.com/golang-grpc/calculator/proto"
	"io"
	"log"
)

func (s *Server) Calculate(ctx context.Context, req *pb.CalulatorRequest) (*pb.CalculatorResponse, error) {
	log.Printf("request values %v,%v,%v\n", req.First, req.Second, req.Operation)

	val := 0

	switch req.Operation {
	case 1:
		val = int(req.First + req.Second)
	case 2:
		val = int(req.First - req.Second)
	case 3:
		val = int(req.First * req.Second)
	default:
		if req.Second != 0 {
			val = int(req.First / req.Second)
		}
	}
	return &pb.CalculatorResponse{
		Result: int32(val),
	}, nil
}

func (s *Server) CalculateNumbers(req *pb.Request, stream pb.CalulatorService_CalculateNumbersServer) error {

	number := req.Number

	k := int32(2)

	for number > 1 {
		if number%k == 0 {
			err := stream.Send(&pb.CalculatorResponse{
				Result: k,
			})
			if err != nil {
				return err
			}
			number /= k
		} else {
			k += 1
		}
	}

	return nil
}

func (s *Server) CalculateAverage(stream pb.CalulatorService_CalculateAverageServer) error {

	totalAmount := 0

	totalCount := 0

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			if totalCount == 0 {
				totalCount = 1
			}
			return stream.SendAndClose(&pb.CalculatorResponse{
				Result: int32(totalAmount / totalCount),
			})
		}

		if err != nil {
			log.Fatalf("error %v\n", err)
		}

		totalAmount += int(msg.Number)
		totalCount++
	}
	return nil

}

func (s *Server) CalculateMax(stream pb.CalulatorService_CalculateMaxServer) error {

	currentMax := -728292929

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if msg.Number > int32(currentMax) {

			currentMax = int(msg.Number)

			err := stream.Send(&pb.CalculatorResponse{
				Result: int32(currentMax),
			})
			if err != nil {
				return err
			}
		}

	}
}
