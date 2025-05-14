FROM golang:1.21 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bff-service ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/bff-service .
COPY --from=builder /app/proto ./proto

EXPOSE 8080
CMD ["./bff-service"]