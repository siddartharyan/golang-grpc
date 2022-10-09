package main

import (
	"context"
	pb "github.com/golang-grpc/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *Server) CreateBlog(ctx context.Context, blog *pb.Blog) (*pb.BlogId, error) {

	blogItem := BlogItem{
		AuthorId: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}

	res, err := collection.InsertOne(ctx, blogItem)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}
	return &pb.BlogId{Id: res.InsertedID.(primitive.ObjectID).Hex()}, nil
}

func (s *Server) ReadBlog(ctx context.Context, id *pb.BlogId) (*pb.Blog, error) {

	oid, err := primitive.ObjectIDFromHex(id.Id)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	data := &BlogItem{}

	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)

	if err := res.Decode(&data); err != nil {
		log.Fatalf("error %v\n", err)
	}

	return documentToBlog(data), nil
}

func (s *Server) UpdateBlog(ctx context.Context, id *pb.Blog) (*emptypb.Empty, error) {

	oid, err := primitive.ObjectIDFromHex(id.Id)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	_, err = collection.UpdateByID(ctx, oid, id)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteBlog(ctx context.Context, id *pb.BlogId) (*emptypb.Empty, error) {

	oid, err := primitive.ObjectIDFromHex(id.Id)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	_, err = collection.DeleteOne(ctx, oid)

	if err != nil {
		log.Fatalf("error %v\n", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListBlogs(_ *emptypb.Empty, stream pb.BlogService_ListBlogsServer) error {

	cur, err := collection.Find(context.Background(), primitive.D{{}})

	if err != nil {
		log.Fatalf("error %v\n", err)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &BlogItem{}
		err := cur.Decode(&data)

		if err != nil {
			log.Fatalf("error %v\n", err)
		}
		stream.Send(documentToBlog(data))
	}
	return nil
}
