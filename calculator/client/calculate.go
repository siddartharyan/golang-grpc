package main

import (
	"context"
	pb "github.com/golang-grpc/calculator/proto"
	"io"
	"log"
)

func doCalculate(c pb.CalulatorServiceClient) {

	res, err := c.Calculate(context.Background(), &pb.CalulatorRequest{
		First:     2,
		Second:    3,
		Operation: 1,
	})

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	log.Println(res.Result)

}

func doCalculateMany(c pb.CalulatorServiceClient) {

	req := &pb.Request{
		Number: 120000,
	}

	stream, err := c.CalculateNumbers(context.Background(), req)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		log.Printf("received number %d\n", msg.Result)
	}

}

func doSendManyTimesReceiveOnce(c pb.CalulatorServiceClient) {

	stream, err := c.CalculateAverage(context.Background())

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	for i := 0; i < 100; i++ {

		err := stream.Send(&pb.Request{
			Number: int32(i),
		})
		if err != nil {
			return
		}

	}

	ans, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	log.Printf("the average value is %d\n\n", ans.Result)

}

func calculateAverage(c pb.CalulatorServiceClient) {

	stream, err := c.CalculateMax(context.Background())

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	numbers := []int{1, 5, 3, 6, 2, 20}

	waitc := make(chan []struct{})

	go func() {
		for _, v := range numbers {
			err := stream.Send(&pb.Request{
				Number: int32(v),
			})
			if err != nil {
				return
			}
		}
		err := stream.CloseSend()
		if err != nil {
			return
		}
	}()

	go func() {

		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error %v\n", err)
			}

			log.Printf("received %v\n", msg.Result)
		}

		close(waitc) //important
	}()

	<-waitc

}
