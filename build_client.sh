#!/bin/bash

# Build for macOS (ARM64)
GOOS=darwin GOARCH=arm64 go build -o cli_chat_client_darwin_arm64 ./cmd/client/main.go

# Build for macOS (AMD64)
GOOS=darwin GOARCH=amd64 go build -o cli_chat_client_darwin_amd64 ./cmd/client/main.go

# Build for Linux (AMD64)
GOOS=linux GOARCH=amd64 go build -o cli_chat_client_linux_amd64 ./cmd/client/main.go

# Build for Windows (AMD64)
export CC=x86_64-w64-mingw32-gcc
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o cli_chat_client_windows_amd64.exe ./cmd/client/main.go