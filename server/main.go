package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/nikfarjam/gRPC-example/pb"
)

func main() {
	addr := ":9292"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error: Unable to run server, %s", err)
	}

	srv := grpc.NewServer()
	var evt Event
	pb.RegisterLogEventServer(srv, &evt)
	reflection.Register(srv)

	log.Printf("Warn: Server is running on %s", addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Error: Server is not able to run, %s", err)
	}
}

func (evt *Event) Start(ctx context.Context, req *pb.StartEventRequest) (*pb.StartEventResponse, error) {
	log.Printf("Info: New start request")
	// TODO: Validate req
	resp := pb.StartEventResponse{
		Guid: req.Guid,
	}

	// TODO: Update statistic
	return &resp, nil
}

func (evt *Event) End(ctx context.Context, req *pb.EndEventRequest) (*pb.EndEventResponse, error) {
	log.Printf("Info: End request")

	resp := pb.EndEventResponse{
		Guid: req.Guid,
	}

	// TODO: Return statistic
	return &resp, nil
}

type Event struct {
	pb.UnimplementedLogEventServer
}
