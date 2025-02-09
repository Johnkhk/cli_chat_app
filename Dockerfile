# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Build the binary with a unique name to avoid conflict with the 'server' directory
RUN CGO_ENABLED=0 GOOS=linux go build -o server-bin ./cmd/server/main.go
# List files to verify that server-bin exists (and note that 'server' is still a directory)
RUN ls -l /app

# Run stage
FROM alpine:latest
WORKDIR /app
# Copy the built binary (server-bin) from the builder stage into the current working directory
COPY --from=builder /app/server-bin .
RUN chmod +x ./server-bin
EXPOSE 50051
CMD ["./server-bin"]
