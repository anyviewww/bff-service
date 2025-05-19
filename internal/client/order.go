package client

import (
	"context"
	"log"

	pb "github.com/anyviewww/bff-service/proto/orders"
	"google.golang.org/grpc"
)

type OrderClient struct {
	client pb.OrderServiceClient
	conn   *grpc.ClientConn
}

func NewOrderClient(conn *grpc.ClientConn) *OrderClient {
	return &OrderClient{
		client: pb.NewOrderServiceClient(conn),
		conn:   conn,
	}
}

func (c *OrderClient) CreateOrder(ctx context.Context, userId uint64, items []int64) (*pb.OrderResponse, error) {
	return c.client.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId: userId,
		Items:  items,
	})
}

func (c *OrderClient) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	return c.client.GetOrder(ctx, req)
}

func (c *OrderClient) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Failed to close order client connection: %v", err)
	}
}
