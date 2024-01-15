# gRPC-example
An example gRPC streaming server with Go Lang

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

### How to run
1. Initiate the module
```bash
go mod init github.com/nikfarjam/gRPC-example
go mod tidy
```
2. Generate gRPC stubs
```bash
rm -r pb
mkdir -p pb
protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative ./logevnt.proto
```
