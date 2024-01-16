package main

import (
	"context"
	"testing"

	"github.com/nikfarjam/gRPC-example/pb"
)

func TestStart(t *testing.T) {
	req := pb.StartEventRequest{
		Guid: "test_1234",
		Name: "Test",
	}

	var srv Event
	resp, err := srv.Start(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Guid != req.Guid {
		t.Fatalf("Guid in response must be request Guid but it was %v", resp.Guid)
	}
}

func TestEnd(t *testing.T) {
	req := pb.EndEventRequest{
		Guid: "test_5678",
	}

	var srv Event
	resp, err := srv.End(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Guid != req.Guid {
		t.Fatalf("Guid in response must be request Guid but it was %v", resp.Guid)
	}
}
