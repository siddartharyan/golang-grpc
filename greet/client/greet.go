package main

import (
	"context"
	"fmt"
	pb "github.com/golang-grpc/greet/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"strconv"
	"time"
)

func doGreet(client pb.GreetServiceClient) {

	res, err := client.Greet(context.Background(), &pb.GreetRequest{
		FirstName: "siddartha",
	})

	if err != nil {

		log.Printf("unable to greet %v\n", err)
	}

	log.Printf("the response from server is %v\n\n", res.GetResult())

}

func doGreetManyTimes(c pb.GreetServiceClient) {

	req := &pb.GreetRequest{
		FirstName: "siddartha",
	}

	stream, err := c.GreetMany(context.Background(), req)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		log.Printf("msg %v\n\n", msg.Result)
	}

}

func doGreetManyTimesandReceiveOnce(c pb.GreetServiceClient) {

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	for i := 0; i < 10; i++ {

		err := stream.Send(&pb.GreetRequest{
			FirstName: fmt.Sprintf("siddartha%d", i),
		})
		if err != nil {
			return
		}

	}

	msg, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	log.Printf("message %v\n\n", msg.Result)

}

func bidirectionalGreet(c pb.GreetServiceClient) {

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	waitc := make(chan []struct{})

	go func() {
		for i := 0; i < 100; i++ {
			err := stream.Send(&pb.GreetRequest{
				FirstName: strconv.FormatInt(int64(1*time.Millisecond), 10),
			})
			if err != nil {
				return
			}
		}
		err := stream.CloseSend() //important
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

func doGreetWithDeadline(c pb.GreetServiceClient, timeout time.Duration) {
	req := pb.GreetRequest{
		FirstName: "siddartha",
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	msg, err := c.GreetWithDeadline(ctx, &req)

	if err != nil {

		e, ok := status.FromError(err)

		if ok {

			if e.Code() == codes.DeadlineExceeded {
				log.Println("deadline exceeded")
				return
			} else {
				log.Fatalf("error is %v\n", e)
			}

		} else {
			log.Fatalf("error is %v\n", err)
		}

	}

	log.Printf("msg is %v\n\n", msg.Result)
}
