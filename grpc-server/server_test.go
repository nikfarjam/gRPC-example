package main

import (
	"context"
	"os"
	"testing"

	"github.com/nikfarjam/gRPC-example/pb"
	"google.golang.org/grpc/metadata"
)

const sec = "test@123"

func TestStart(t *testing.T) {
	os.Setenv(TOKEN, sec)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{TOKEN: []string{sec}})

	req := pb.StartEventRequest{
		Guid: "test_1234",
		Name: "Test",
	}

	var srv Server
	resp, err := srv.Start(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Guid != req.Guid {
		t.Fatalf("Guid in response must be request Guid but it was %v", resp.Guid)
	}
}

func TestEnd(t *testing.T) {
	os.Setenv(TOKEN, sec)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{TOKEN: []string{sec}})

	req := pb.EndEventRequest{
		Guid: "test_5678",
	}

	var srv Server
	resp, err := srv.End(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Guid != req.Guid {
		t.Fatalf("Guid in response must be request Guid but it was %v", resp.Guid)
	}
}
