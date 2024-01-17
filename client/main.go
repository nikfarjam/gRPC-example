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
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nikfarjam/gRPC-example/pb"
)

var URLS = [...]string{"/home/", "/faq/", "/login", "/404", "/asset.js", "/main.css", "/intranet-analytics/", "index.html"}
var METHODS = [...]string{"GET", "POST"}

const TOKEN = "log_server_token"

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
	ctx = metadata.AppendToOutgoingContext(ctx, TOKEN, os.Getenv(TOKEN))
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
	ctx = metadata.AppendToOutgoingContext(ctx, TOKEN, os.Getenv(TOKEN))
	defer cancel()

	resp, err := c.Start(ctx, &req)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Println(resp)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, TOKEN, os.Getenv(TOKEN))
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

	var delay int64
	if len(os.Args) > 3 {
		delay, _ = strconv.ParseInt(os.Args[3], 10, 64)
	}

	for i := 0; i < numMessages; i += 1 {
		time.Sleep(time.Duration(delay) * time.Microsecond)
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
			log.Fatalf("error: %s", err)
		}
	}
	sresp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Println(sresp)

	ereq := pb.EndEventRequest{
		Guid: guid,
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, TOKEN, os.Getenv(TOKEN))
	defer cancel()
	eresp, err := c.End(ctx, &ereq)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Printf("Total result is %+v\n", eresp)
}
