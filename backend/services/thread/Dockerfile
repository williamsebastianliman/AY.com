# builder stage
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app

# Copy everything so Go can see your go.work / multi‐module layout
COPY . .

# Build the thread binary (uses your go.work & generated protos)
RUN cd services/thread && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /thread ./cmd/thread

# final runtime image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the built binary
COPY --from=builder /thread .

# Create a non-root user and switch
RUN addgroup -S app && adduser -S -G app app
USER app

EXPOSE 50051
ENTRYPOINT ["./thread"]