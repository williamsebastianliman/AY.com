FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY services/media/go.mod services/media/go.sum ./services/media/
COPY proto/gen/go ./proto/gen/go
COPY services/media ./services/media
RUN cd services/media && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/media ./cmd/media

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/media .
RUN addgroup -S app && adduser -S -G app app
USER app
EXPOSE 50051
ENTRYPOINT ["./media"]