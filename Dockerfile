# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go
RUN ls -l /app  # Verify that the binary 'server' exists

# Run stage
FROM alpine:latest
WORKDIR /app
# Copy the built binary into the current working directory (i.e. /app)
COPY --from=builder /app/server .
RUN chmod +x ./server
EXPOSE 50051
CMD ["./server"]
