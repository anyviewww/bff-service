# bff-service
protoc --go_out=. --go-grpc_out=. proto/menu/menu.proto
protoc --go_out=. --go-grpc_out=. proto/order/order.proto

`` go run cmd/server/main.go ``
