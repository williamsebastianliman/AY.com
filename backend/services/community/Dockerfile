FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY services/community/go.mod services/community/go.sum ./services/community/
COPY proto/gen/go ./proto/gen/go
COPY services/community ./services/community
RUN cd services/community && go mod download && go build -o /app/community ./cmd/community

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/community .
EXPOSE 50051
CMD ["./community"]