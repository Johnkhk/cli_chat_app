# Generate gRPC code
```proto
find proto -name "*.proto" | xargs protoc --go_out=genproto --go-grpc_out=genproto --go_opt=module=github.com/johnkhk/cli_chat_app/genproto --go-grpc_opt=module=github.com/johnkhk/cli_chat_app/genproto
```