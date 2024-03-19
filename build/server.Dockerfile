# Build the binary and run CI
FROM golang:1.22-alpine AS builder

# Get certs
RUN apk --update add ca-certificates

# Copy local source
WORKDIR /app
COPY . .

# Build binary
ARG KEY
RUN GOOS=linux go build -a -o server cmd/server/main.go

# Build image
FROM scratch

# Copy certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy config and binary
ARG KEY
WORKDIR /root/
COPY --from=builder /app/server .

# Entry and port
CMD ["./server"]
EXPOSE 8080
