FROM golang:1.21.4-alpine3.18 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bff-service ./cmd/server

FROM alpine:3.18.4
WORKDIR /app
COPY --from=builder /app/bff-service /app/
COPY --from=builder /app/proto/ /app/proto/
COPY --from=builder /app/configs/ /app/configs/
EXPOSE 8080
CMD ["./bff-service"]