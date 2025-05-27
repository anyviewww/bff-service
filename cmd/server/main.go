package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/anyviewww/bff-service/internal/api"
	"github.com/anyviewww/bff-service/internal/config"
	pbDishes "github.com/anyviewww/bff-service/proto/dishes"
	pbOrders "github.com/anyviewww/bff-service/proto/orders"
)

func main() {
	cfg := config.Load()

	// initialising gRPC connections
	menuConn := createGRPCConnection(cfg.MenuServiceAddr)
	defer menuConn.Close()

	orderConn := createGRPCConnection(cfg.OrderServiceAddr)
	defer orderConn.Close()

	// customer creation
	dishClient := pbDishes.NewDishServiceClient(menuConn)
	orderClient := pbOrders.NewOrderServiceClient(orderConn)

	// configuring the HTTP server
	router := gin.Default()

	// create API handler
	apiHandler := api.NewHandler(dishClient, orderClient, cfg)
	apiRouter := api.NewRouter(apiHandler)
	apiRouter.SetupRoutes(router)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server started on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func createGRPCConnection(addr string) *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC service at %s: %v", addr, err)
	}

	return conn
}
