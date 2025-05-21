package client

import (
	"context"
	"log"

	pb "github.com/anyviewww/bff-service/proto/dishes"
	"google.golang.org/grpc"
)

type MenuClient struct {
	client pb.DishServiceClient
	conn   *grpc.ClientConn
}

// GetDishes implements dishes.DishServiceClient.
func (c *MenuClient) GetDishes(ctx context.Context, in *pb.DishRequest, opts ...grpc.CallOption) (*pb.DishesResponse, error) {
	panic("unimplemented")
}

func NewMenuClient(conn *grpc.ClientConn) *MenuClient {
	return &MenuClient{
		client: pb.NewDishServiceClient(conn),
		conn:   conn,
	}
}

func (c *MenuClient) GetDish(ctx context.Context, id int32) (*pb.Dish, error) {
	resp, err := c.client.GetDishes(ctx, &pb.DishRequest{Id: id})
	if err != nil {
		return nil, err
	}
	if len(resp.Dishes) == 0 {
		return nil, nil
	}
	return resp.Dishes[0], nil
}

func (c *MenuClient) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Failed to close menu connection: %v", err)
	}
}
