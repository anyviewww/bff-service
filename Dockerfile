FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o bff-service ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/bff-service .
COPY --from=builder /app/proto ./proto
EXPOSE 8080
CMD ["./bff-service"]