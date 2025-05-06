package api

import (
	"bff-service/grpc/orders"
	pb "bff-service/grpc/orders"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	UserID uint64  `json:"user_id"`
	Items  []int64 `json:"items"`
}

type GetOrderRequest struct {
	ID uint64 `json:"id"`
}

type UpdateOrderRequest struct {
	ID     uint64  `json:"id"`
	UserID uint64  `json:"user_id"`
	Items  []int64 `json:"items"`
	Status string  `json:"status"`
}

type DeleteOrderRequest struct {
	ID uint64 `json:"id"`
}

var baseHandler = BaseHandler{}

func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if !baseHandler.BindAndValidate(c, &req) {
		return
	}

	ctx := context.Background()
	resp, err := orders.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId: req.UserID,
		Items:  req.Items,
	})
	baseHandler.HandleGRPCError(c, err)
	baseHandler.Respond(c, http.StatusOK, gin.H{
		"id":      resp.Id,
		"user_id": resp.UserId,
		"items":   resp.Items,
		"status":  resp.Status,
	})
}

func GetOrder(c *gin.Context) {
	var req GetOrderRequest
	if !baseHandler.BindAndValidate(c, &req) {
		return
	}

	ctx := context.Background()
	resp, err := orders.GetOrder(ctx, &pb.GetOrderRequest{
		Id: req.ID,
	})
	baseHandler.HandleGRPCError(c, err)
	baseHandler.Respond(c, http.StatusOK, gin.H{
		"id":      resp.Id,
		"user_id": resp.UserId,
		"items":   resp.Items,
		"status":  resp.Status,
	})
}

func UpdateOrder(c *gin.Context) {
	var req UpdateOrderRequest
	if !baseHandler.BindAndValidate(c, &req) {
		return
	}

	ctx := context.Background()
	resp, err := orders.UpdateOrder(ctx, &pb.UpdateOrderRequest{
		Id:     req.ID,
		UserId: req.UserID,
		Items:  req.Items,
		Status: req.Status,
	})
	baseHandler.HandleGRPCError(c, err)
	baseHandler.Respond(c, http.StatusOK, gin.H{
		"id":      resp.Id,
		"user_id": resp.UserId,
		"items":   resp.Items,
		"status":  resp.Status,
	})
}

func DeleteOrder(c *gin.Context) {
	var req DeleteOrderRequest
	if !baseHandler.BindAndValidate(c, &req) {
		return
	}

	ctx := context.Background()
	resp, err := orders.DeleteOrder(ctx, &pb.DeleteOrderRequest{
		Id: req.ID,
	})
	baseHandler.HandleGRPCError(c, err)
	baseHandler.Respond(c, http.StatusOK, gin.H{"deleted": resp.Deleted})
}
