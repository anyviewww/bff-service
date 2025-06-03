package api

import (
	"github.com/anyviewww/bff-service/internal/config"
	pbDishes "github.com/anyviewww/bff-service/proto/dishes"
	pbOrders "github.com/anyviewww/bff-service/proto/orders"
)

type Handler struct {
	menuClient  pbDishes.DishServiceClient
	orderClient pbOrders.OrderServiceClient
	cfg         *config.Config
}

func NewHandler(menuClient pbDishes.DishServiceClient, orderClient pbOrders.OrderServiceClient, cfg *config.Config) *Handler {
	return &Handler{
		menuClient:  menuClient,
		orderClient: orderClient,
		cfg:         cfg,
	}
}
