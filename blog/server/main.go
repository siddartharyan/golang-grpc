package main

import (
	"context"
	pb "github.com/golang-grpc/blog/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	pb.BlogServiceServer
}

var address = "0.0.0.0:8080"

var collection *mongo.Collection

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	err = client.Connect(context.Background())

	collection = client.Database("blogDb").Collection("blog")

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	s := grpc.NewServer()

	pb.RegisterBlogServiceServer(s, &Server{})

	log.Printf("server is started %v\n\n", address)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("error %v\n", err)
	}

}
