package orders

import (
	pb "bff-service/grpc/orders"
	"context"
	"log"

	"google.golang.org/grpc"
)

var orderServiceClient pb.OrderServiceClient

func InitGRPCConnections() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to OrderService: %v", err)
	}
	orderServiceClient = pb.NewOrderServiceClient(conn)
}

func CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	return orderServiceClient.CreateOrder(ctx, req)
}

func GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	return orderServiceClient.GetOrder(ctx, req)
}

func UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
	return orderServiceClient.UpdateOrder(ctx, req)
}

func DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	return orderServiceClient.DeleteOrder(ctx, req)
}
