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
# Copy the binary from the builder stage explicitly to /app/server
COPY --from=builder /app/server /app/server
# Ensure the binary is executable
RUN chmod +x /app/server
EXPOSE 50051
# Use the full path to run the server binary
CMD ["/app/server"]
