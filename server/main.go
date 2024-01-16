package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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
	log.Printf("Info: Hand shake with new client [%v]", req.Name)
	// TODO: Validate req
	resp := pb.StartEventResponse{
		Guid: req.Guid,
	}

	// TODO: Update statistic
	return &resp, nil
}

func (evt *Event) Event(stream pb.LogEvent_EventServer) error {
	guid := ""
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Can not read the stream")
		}
		log.Printf("Info: New event from client [%v]", req.Guid)
		log.Printf("Debug: New event from client [%v]", req)
		guid = req.Guid
	}

	resp := pb.LogEventResponse{
		Guid: guid,
	}

	return stream.SendAndClose(&resp)
}

func (evt *Event) End(ctx context.Context, req *pb.EndEventRequest) (*pb.EndEventResponse, error) {
	log.Printf("Info: End with client [%v]", req.Guid)

	resp := pb.EndEventResponse{
		Guid: req.Guid,
	}

	// TODO: Return statistic
	return &resp, nil
}

type Event struct {
	pb.UnimplementedLogEventServer
}
