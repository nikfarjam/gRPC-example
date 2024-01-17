package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/nikfarjam/gRPC-example/pb"
)

const TOKEN = "log_server_token"

func main() {
	addr := ":9292"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error: Unable to run server, %s", err)
	}

	srv := createServer()

	log.Printf("Warn: Server is running on %s", addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Error: Server is not able to run, %s", err)
	}
}

func createServer() *grpc.Server {
	server := grpc.NewServer(grpc.UnaryInterceptor(timingInterceptor))
	var u Server
	pb.RegisterLogEventServer(server, &u)
	reflection.Register(server)
	return server
}

func timingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		log.Printf("info: %s took %v", info.FullMethod, duration)
	}()

	return handler(ctx, req)
}

func (srv *Server) Start(ctx context.Context, req *pb.StartEventRequest) (*pb.StartEventResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata")
	}
	if len(md[TOKEN]) == 0 || md[TOKEN][0] != os.Getenv(TOKEN) {
		return nil, status.Errorf(codes.Unauthenticated, "Token is not valid")
	}

	log.Printf("Info: Hand shake with new client [%v]", req.Name)
	if req.Guid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Guid is empty")
	}
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name is empty")
	}

	resp := pb.StartEventResponse{
		Guid: req.Guid,
	}

	if srv.Statistic == nil {
		srv.Statistic = make(map[string]Statistic)
	}
	srv.Statistic[req.Guid] = *NewStatistic()

	return &resp, nil
}

func (srv *Server) Event(stream pb.LogEvent_EventServer) error {
	guid := ""
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Can not read the stream")
		}
		log.Printf("Debug: New event from client [%v]", req.Guid)
		guid = req.Guid
		s := srv.Statistic[req.Guid]
		s.Ips[req.Ip] = s.Ips[req.Ip] + 1
		s.Urls[req.Path] = s.Urls[req.Path] + 1
	}

	resp := pb.LogEventResponse{
		Guid: guid,
	}

	return stream.SendAndClose(&resp)
}

func (srv *Server) End(ctx context.Context, req *pb.EndEventRequest) (*pb.EndEventResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata")
	}
	if len(md[TOKEN]) == 0 || md[TOKEN][0] != os.Getenv(TOKEN) {
		return nil, status.Errorf(codes.Unauthenticated, "Token is not valid")
	}
	if req.Guid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Guid is empty")
	}
	s := srv.Statistic[req.Guid]
	topClientIps := topItemsByValue(s.Ips, 3)
	topVisitedUrls := topItemsByValue(s.Urls, 3)

	resp := pb.EndEventResponse{
		Guid:           req.Guid,
		NumUniqueIp:    int64(len(s.Ips)),
		TopClientIps:   topClientIps,
		TopVisitedUrls: topVisitedUrls,
	}
	log.Printf("Debug: End with client %v", resp.Guid)

	return &resp, nil
}

func topItemsByValue[T string | float32 | int64, P float32 | int64 | int](m map[T]P, s int) []T {
	keys := make([]T, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] < m[keys[j]]
	})

	if len(keys) < s {
		return keys
	}
	return keys[:s]
}

type Server struct {
	Statistic map[string]Statistic
	pb.UnimplementedLogEventServer
}

type Statistic struct {
	Urls map[string]int64
	Ips  map[string]int64
}

func NewStatistic() *Statistic {
	s := Statistic{
		Urls: make(map[string]int64),
		Ips:  make(map[string]int64),
	}
	return &s
}
