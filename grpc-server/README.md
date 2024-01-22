# gRPC Server
An example Observability Server that captures and analyze applications' logs via the gRPC protocol. This server actively listens to incoming log streams from clients and returns statistics of tops visited URLS, top most common IPS and number of unique IPS on demand.

This server has one end point with three APIs. For authorization, Server checks if if both Server and Client has the same token. gRPCurl can be used to get the end point details. 

## How to install
### Prerequisites
1. [Go Compiler](https://go.dev/doc/install)
2. [Protocol Buffers compiler](https://grpc.io/docs/protoc-installation/)
3. [GoLand](https://www.jetbrains.com/go/) or [Visual Studio Code](https://code.visualstudio.com/download)
4. **Go plugins** for the protocol compiler that generates required Go code from *.proto* files
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
export PATH="$PATH:$(go env GOPATH)/bin"
```
5. [Postman](https://www.postman.com/downloads/) or [grpCurl](https://github.com/fullstorydev/grpcurl/releases) to test the server

### How to run
1. Initiate the module
```bash
git clone https://github.com/nikfarjam/gRPC-example.git
cd ./grpc-server
go mod init
go mod tidy
```
2. Generate gRPC stubs
```bash
rm -r pb && mkdir pb
protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative ./logevnt.proto
```

### Run Server
```bash
# note: log_server_token must be lowercase
export log_server_token=<secret-token>
go run .
```

## API Guide

### Endpoint description

To assist Clients at runtime to construct RPC requests and responses and also for documentation the server use reflection. Example commands to query EndPoint description

```bash
grpcurl --plaintext localhost:9292 list LogEvent
grpcurl --plaintext localhost:9292 describe LogEvent.Start
```

### APIs
1. **Start:** Register a client with a name and unique Guid
- **Request**  

Example body
```
{
    "guid": "428fd17a3c92f1ebac3ac73f1ccaca2f1abfc5b4",
    "name": "webServer"
}
```
Example request
```bash
export log_server_token=<secret-token>
grpcurl --plaintext -d @ -H 'log_server_token:'"$log_server_token"''  localhost:9292 LogEvent.Start <  ./examples/start_req.json
```
- **Response**  

Example body
```
{
    "guid": "428fd17a3c92f1ebac3ac73f1ccaca2f1abfc5b4"
}
```

2. **Event:** A stream API which listen to client Log requests that contains visited IP, Path and store them in memory
- Request  
Example body
```
{
    "guid": "cc52a3c9b23ec23349a0310d4770ef000b93b08b",
    "ip": "10.0.1.12",
    "time": "2020-07-30T18:00:00.000Z",
    "method": "Get",
    "path": "/home/",
    "status": 200,
    "process_time": 231,
    "full_log": "10.0.1.12 - - [07/Jan/2024:16:01:17 +0200] \"GET /home/ HTTP/1.1\" 200 231 \"-\" \"Mozilla/5.0 (X11; U; Linux x86_64; fr-FR) AppleWebKit/534.7 (KHTML, like Gecko) Epiphany/2.30.6 Safari/534.7\""
}
```
Example request
```bash
grpcurl --plaintext -d @ -H 'log_server_token:'"$log_server_token"''  localhost:9292 LogEvent.Event <  ./examples/event_req.json
```
- Response  
Example body
```
{
    "guid": "cc52a3c9b23ec23349a0310d4770ef000b93b08b"
}
```

3. **End:** Return number of unique IPs, most visited Path and IPs
- **Request**  

Example body
```
{
    "guid": "428fd17a3c92f1ebac3ac73f1ccaca2f1abfc5b4",
}
```
Example request
```bash
grpcurl --plaintext -d @ -H 'log_server_token:'"$log_server_token"''  localhost:9292 LogEvent.End <  ./examples/end_event_req.json
```
- **Response**  

Example body
```
{
  "guid": "428fd17a3c92f1ebac3ac73f1ccaca2f1abfc5b4",
  "numUniqueIp": "1",
  "topClientIps": [
    "10.0.1.12"
  ],
  "topVisitedUrls": [
    "/home/"
  ]
}
```

### Test
#### With Go test
`go test -v ./...`

#### with gRPCurl
```bash
# Set token
export log_server_token=<secret-token>
# Register client
grpcurl --plaintext -d @ -H 'log_server_token:'"$log_server_token"''  localhost:9292 LogEvent.Start <  ./examples/start_req.json
# Send log stream
grpcurl --plaintext -d @ -H 'log_server_token:'"$log_server_token"''  localhost:9292 LogEvent.Event <  ./examples/event_req.json
# Send get statistic request
grpcurl --plaintext -d @ -H 'log_server_token:'"$log_server_token"''  localhost:9292 LogEvent.End <  ./examples/end_event_req.json
```

#### Mock Client
```bash
export log_server_token=<secret-token>
go run ../grpc-client/ <client name> <number of messages> <delay between each messages>
// For example 
// $ go run ./client test-1 100 5
```
**Note:** Secret token for both server and client terminal must be the same 