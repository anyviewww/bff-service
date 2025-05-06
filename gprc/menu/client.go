package menu

import (
	pb "bff-service/grpc/menu"
	"context"
	"log"

	"google.golang.org/grpc"
)

var dishServiceClient pb.DishServiceClient

func InitGRPCConnections() {
	conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to DishService: %v", err)
	}
	dishServiceClient = pb.NewDishServiceClient(conn)
}

func GetDishes(ctx context.Context, req *pb.DishRequest) (*pb.DishesResponse, error) {
	return dishServiceClient.GetDishes(ctx, req)
}
