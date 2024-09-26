# generate proto/auth.proto

protoc --go_out=./genproto --go_opt=module=github.com/johnkhk/cli_chat_app/proto \
  --go-grpc_out=./genproto --go-grpc_opt=module=github.com/johnkhk/cli_chat_app/proto \
  proto/auth/auth.proto

# generate proto/friends.proto

protoc --go_out=./genproto --go_opt=module=github.com/johnkhk/cli_chat_app/proto \
  --go-grpc_out=./genproto --go-grpc_opt=module=github.com/johnkhk/cli_chat_app/proto \
  proto/friends/friends.proto


# generate proto/chat.proto

protoc --go_out=./genproto --go_opt=module=github.com/johnkhk/cli_chat_app/proto \
  --go-grpc_out=./genproto --go-grpc_opt=module=github.com/johnkhk/cli_chat_app/proto \
  proto/chat/chat.proto
