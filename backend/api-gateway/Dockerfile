FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.work go.work.sum ./            
COPY api-gateway/go.mod api-gateway/go.sum ./

COPY . .

WORKDIR /app/api-gateway
RUN go mod tidy
RUN go build -o /usr/local/bin/api-gateway ./cmd/gateway

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /usr/local/bin/api-gateway .
EXPOSE 8080
ENTRYPOINT ["./api-gateway"]