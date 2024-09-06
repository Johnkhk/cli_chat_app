# generate proto/auth.proto

protoc --go_out=./genproto --go_opt=module=github.com/johnkhk/cli_chat_app/proto \
  --go-grpc_out=./genproto --go-grpc_opt=module=github.com/johnkhk/cli_chat_app/proto \
  proto/auth/auth.proto


protoc --go_out=./genproto --go_opt=module=github.com/johnkhk/cli_chat_app/proto \
  --go-grpc_out=./genproto --go-grpc_opt=module=github.com/johnkhk/cli_chat_app/proto \
  proto/message/message.proto