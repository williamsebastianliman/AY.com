FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY services/auth/go.mod services/auth/go.sum ./services/auth/
COPY proto/gen/go ./proto/gen/go
COPY services/auth ./services/auth
RUN cd services/auth && go mod download && go build -o /app/auth ./cmd/auth

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth .
EXPOSE 50051
CMD ["./auth"]