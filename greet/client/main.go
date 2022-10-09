package main

import (
	pb "github.com/golang-grpc/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

var address = "localhost:50051"

func main() {

	tls := true

	opts := []grpc.DialOption{}

	if tls {

		certFile := "../ssl/ca.crt"

		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("error is %v\n", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	// to avoid ssl in local
	//conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(address, opts...)

	if err != nil {

		log.Fatalf("error in connecting to client %v\n", err)
	}

	c := pb.NewGreetServiceClient(conn)

	doGreet(c)
	doGreetManyTimes(c)
	doGreetManyTimesandReceiveOnce(c)
	bidirectionalGreet(c)
	doGreetWithDeadline(c, 4*time.Second)
	defer conn.Close()

}
