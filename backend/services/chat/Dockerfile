FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY services/chat/go.mod services/chat/go.sum ./services/chat/
COPY proto/gen/go ./proto/gen/go
COPY services/chat ./services/chat
RUN cd services/chat && go mod download && go build -o /app/chat ./cmd/chat

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/chat .
EXPOSE 50051
CMD ["./chat"]