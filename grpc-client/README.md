# gRPC Server
An example gRPC client to test [Observability Server](../grpc-server/)

#### Mock Client
```bash
export log_server_token=<secret-token>
go run ../grpc-server &
go run . <client name> <number of messages> <delay between each messages>
// For example 
// $ go run ./client test-1 100 5
```