# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go
RUN ls -l /app # List all files in /app to verify that server exists

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 50051
CMD ["./main"]  # Run the 'main' binary
