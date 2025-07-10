# Start from a small Go base image to build the app
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.work and go.work.sum and monorepo modules is fast solution but it's not recommended
# Copy go.work and go.work.sum to leverage Docker layer caching
COPY go.work go.work.sum ./

# Copy go.mod and go.sum first to leverage Docker layer caching

COPY telemetry-node/go.mod telemetry-node/go.sum ./telemetry-node/
COPY telemetry-sink/go.mod telemetry-sink/go.sum ./telemetry-sink/
COPY pkg/go.mod pkg/go.sum ./pkg/

#WORKDIR /app/telemetry-node

RUN go mod download

# Copy the source code
COPY telemetry-node/ ./telemetry-node/
COPY telemetry-sink/ ./telemetry-sink/
COPY pkg/ ./pkg/

# Build the Go binary
RUN go build -o telemetry-sink-app ./telemetry-node/cmd

# --- Runtime image ---
FROM alpine:latest

# Optional: Install ca-certificates if your app makes HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/telemetry-sink-app .

# Run the app
ENTRYPOINT ["/app/telemetry-sink-app"]
