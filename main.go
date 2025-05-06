package main

import (
	"bff-service/api"
	"bff-service/grpc/menu"
	"bff-service/grpc/orders"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initializing gRPC clients
	orders.InitGRPCConnections()
	menu.InitGRPCConnections()

	// Create a new Gin router
	r := gin.Default()

	// Route group for AuthService
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", api.Login)
		authGroup.POST("/register", api.Register)
		authGroup.POST("/reset-password", api.ResetPassword)
	}

	// Protected Route Group
	protected := r.Group("/")
	protected.Use(api.AuthMiddleware())
	{
		// Route group for OrderService
		orderGroup := protected.Group("/order")
		{
			orderGroup.POST("/create", api.CreateOrder)
			orderGroup.POST("/get", api.GetOrder)
			orderGroup.PUT("/update", api.UpdateOrder)
			orderGroup.DELETE("/delete", api.DeleteOrder)
		}

		// Route group for DishService
		menuGroup := protected.Group("/menu")
		{
			menuGroup.POST("/get-dishes", api.GetDishes)
		}
	}

	// We start the server on port 8080
	r.Run(":8080")
}
