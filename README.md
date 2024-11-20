## Building the CLI with a Custom Server Address

To build with a default server address:
```bash
go build -ldflags="-X 'main.defaultServerAddress=<your-server-ip>'" -o cli_chat_app
```

To override the server address at runtime:
```bash
SERVER_ADDRESS=<your-server-ip> ./cli_chat_app
```