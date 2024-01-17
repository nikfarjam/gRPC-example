# gRPC-example
An example gRPC streaming server with Go Lang  
This server has three API  
1. **Start:** Register a client with a name and unique Guid
2. **Event:** A stream API which listen to client Log requests that contains visited IP, Path and store them in memory
3. **End:** Return number of unique IPs, most visited Path and IPs

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
go mod init github.com/nikfarjam/gRPC-example
go mod tidy
```
2. Generate gRPC stubs
```bash
rm -r pb && mkdir pb
protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative ./logevnt.proto
```

### Run Server
```bash
export LOG_SERVER_TOKEN=<secret-token>
go run ./server/
```

### Test server
#### With Go test
`go test -v ./...`

#### with gRPCurl
```bash
grpcurl --plaintext localhost:9292 list LogEvent
grpcurl --plaintext localhost:9292 describe LogEvent.Start

grpcurl --plaintext -d @ localhost:9292 LogEvent.Start <  ./examples/start_req.json
```

#### Mock Client
```bash
export LOG_SERVER_TOKEN=<secret-token>
go run ./client/ <name> <number of messages> <delay between each messages>
// For example 
// $ go run ./client test-1 100 5
```
**Note:** Secret token for both server and client terminal must be the same 