package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/nikfarjam/gRPC-example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var URLS = [...]string{"/home/", "/faq/", "/login", "/404", "/asset.js", "/main.css", "/intranet-analytics/", "index.html"}
var METHODS = [...]string{"GET", "POST"}

func TestStartE2E(t *testing.T) {
	os.Setenv(TOKEN, sec)
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	srv := createServer()
	go srv.Serve(lis)

	port := lis.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("localhost:%d", port)
	creds := insecure.NewCredentials()
	conn, err := grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewLogEventClient(conn)

	req := pb.StartEventRequest{
		Guid: "test_1234",
		Name: "Test",
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), TOKEN, os.Getenv(TOKEN))
	resp, err := c.Start(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Guid != req.Guid {
		t.Fatalf("bad response Guid: got %#v, expected %#v", resp.Guid, req.Guid)
	}
}

func TestEndE2E(t *testing.T) {
	os.Setenv(TOKEN, sec)
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	srv := createServer()
	go srv.Serve(lis)

	port := lis.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("localhost:%d", port)
	creds := insecure.NewCredentials()
	conn, err := grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewLogEventClient(conn)

	req := pb.EndEventRequest{
		Guid: "test_1234",
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), TOKEN, os.Getenv(TOKEN))
	resp, err := c.End(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Guid != req.Guid {
		t.Fatalf("bad response Guid: got %#v, expected %#v", resp.Guid, req.Guid)
	}
}

func TestEventE2E(t *testing.T) {
	guid := "test_1234"
	numMessages := 5
	os.Setenv(TOKEN, sec)
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	srv := createServer()
	go srv.Serve(lis)

	port := lis.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("localhost:%d", port)
	creds := insecure.NewCredentials()
	conn, err := grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewLogEventClient(conn)

	req := pb.StartEventRequest{
		Guid: guid,
		Name: "Test",
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), TOKEN, os.Getenv(TOKEN))
	_, err = c.Start(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}
	stream, err := c.Event(ctx)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numMessages; i += 1 {
		lreq := pb.LogEventRequest{
			Guid:        guid,
			Ip:          fmt.Sprintf("%v.%v.%v.%v", r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255)),
			Time:        timestamppb.Now(),
			Method:      METHODS[r.Intn(len(METHODS))],
			Path:        URLS[r.Intn(len(URLS))],
			Status:      r.Int31n(5)*100 + r.Int31n(10)*9,
			ProcessTime: r.Int63n(800),
			FullLog:     fmt.Sprintf("Dumb log %v", i+1),
		}
		if err := stream.Send(&lreq); err != nil {
			t.Fatalf("error: %s", err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	ereq := pb.EndEventRequest{
		Guid: guid,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, TOKEN, os.Getenv(TOKEN))
	defer cancel()
	eresp, err := c.End(ctx, &ereq)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if eresp.Guid != guid {
		t.Fatalf("bad response Guid: got %#v, expected %#v", eresp.Guid, guid)
	}
	if int(eresp.NumUniqueIp) != numMessages {
		t.Fatalf("bad response NumUniqueIp: got %#v, expected %#v", eresp.NumUniqueIp, numMessages)
	}
	if len(eresp.TopClientIps) != 3 {
		t.Fatalf("bad response TopClientIps has %#v items but must have %#v", eresp.NumUniqueIp, 3)
	}
	if len(eresp.GetTopVisitedUrls()) != 3 {
		t.Fatalf("bad response GetTopVisitedUrls has %#v items but must have %#v", eresp.NumUniqueIp, 3)
	}

}
