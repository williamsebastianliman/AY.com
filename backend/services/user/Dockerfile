FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY services/user/go.mod services/user/go.sum ./services/user/
COPY proto/gen/go ./proto/gen/go
COPY services/user ./services/user
RUN cd services/user && go mod download && go build -o /app/user ./cmd/user

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user .
EXPOSE 50051
CMD ["./user"]