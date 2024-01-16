package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nikfarjam/gRPC-example/pb"
)

var IPS = [...]string{"127.0.0.1", "10.0.0.12", "8.8.8.8", "10.0.1.20"}
var URLS = [...]string{"/home/", "/faq/", "/login", "/404"}
var METHODS = [...]string{"GET", "POST"}

func main() {
	var name string
	if len(os.Args) < 2 {
		name = "Go Client"
	} else {
		name = os.Args[1]
	}
	h := sha256.New()
	h.Write([]byte(name))

	guid := fmt.Sprintf("%x", h.Sum(nil)[:20])

	fmt.Printf("Client Name is '%v'\n", name)

	addr := "localhost:9292"
	creds := insecure.NewCredentials()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer conn.Close()

	log.Printf("Info: connected to %s", addr)
	c := pb.NewLogEventClient(conn)

	req := pb.StartEventRequest{
		Guid: guid,
		Name: name,
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.Start(ctx, &req)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Println(resp)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.Event(ctx)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numMessages := -1
	if len(os.Args) > 2 {
		numMessages, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Can't convert %v to an int!", os.Args[2])
		}
	}
	if numMessages < 0 {
		numMessages = 3
	}

	for i := 0; i < numMessages; i += 1 {
		lreq := pb.LogEventRequest{
			Guid:        guid,
			Ip:          IPS[r.Intn(len(IPS))],
			Time:        timestamppb.Now(),
			Method:      METHODS[r.Intn(len(METHODS))],
			Path:        URLS[r.Intn(len(URLS))],
			Status:      r.Int31n(5)*100 + r.Int31n(10)*9,
			ProcessTime: r.Int63n(800),
			FullLog:     "Dumb log",
		}
		if err := stream.Send(&lreq); err != nil {
			log.Fatalf("error: %s", err)
		}
	}
	lresp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Println(lresp)

	ereq := pb.EndEventRequest{
		Guid: guid,
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	eresp, err := c.End(ctx, &ereq)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Println(eresp)
}
