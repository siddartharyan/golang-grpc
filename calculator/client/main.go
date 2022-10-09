package main

import (
	pb "github.com/golang-grpc/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var address = "localhost:8081"

func main() {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	c := pb.NewCalulatorServiceClient(conn)

	doCalculate(c)

	doCalculateMany(c)

	doSendManyTimesReceiveOnce(c)

	calculateAverage(c)

	defer conn.Close()

}
