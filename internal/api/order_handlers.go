package api

import (
	"net/http"
	"strconv"

	pbOrders "github.com/anyviewww/bff-service/proto/orders"
	"github.com/gin-gonic/gin"
)

type orderRequest struct {
	Items  []int64 `json:"items,omitempty"`
	Status *string `json:"status,omitempty"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req struct {
		Items []int64 `json:"items" binding:"required,min=1"`
	}
	if !h.bindJSON(c, &req) {
		return
	}

	order, ok := h.createOrder(c, userID, req.Items)
	if !ok {
		return
	}

	c.JSON(http.StatusCreated, h.toOrderResponse(order))
}

func (h *Handler) GetOrder(c *gin.Context) {
	order, ok := h.getValidatedOrder(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, h.toOrderResponse(order))
}

func (h *Handler) UpdateOrder(c *gin.Context) {
	order, ok := h.getValidatedOrder(c)
	if !ok {
		return
	}

	var req orderRequest
	if !h.bindJSON(c, &req) {
		return
	}

	updatedOrder, ok := h.updateOrder(c, order.Id, req)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, h.toOrderResponse(updatedOrder))
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	order, ok := h.getValidatedOrder(c)
	if !ok {
		return
	}

	if !h.deleteOrder(c, order.Id) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func (h *Handler) GetUserOrders(c *gin.Context) {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.orderClient.GetUserOrders(c.Request.Context(), &pbOrders.GetUserOrdersRequest{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]gin.H, 0, len(orders.Orders))
	for _, order := range orders.Orders {
		response = append(response, h.toOrderResponse(order))
	}

	c.JSON(http.StatusOK, gin.H{"orders": response})
}

func (h *Handler) createOrder(c *gin.Context, userID uint64, items []int64) (*pbOrders.OrderResponse, bool) {
	order, err := h.orderClient.CreateOrder(c.Request.Context(), &pbOrders.CreateOrderRequest{
		UserId: userID,
		Items:  items,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}
	return order, true
}

func (h *Handler) updateOrder(c *gin.Context, orderID uint64, req orderRequest) (*pbOrders.OrderResponse, bool) {
	updateReq := &pbOrders.UpdateOrderRequest{
		Id: orderID,
	}

	if req.Items != nil {
		updateReq.Items = req.Items
	}
	if req.Status != nil {
		updateReq.Status = *req.Status
	}

	updatedOrder, err := h.orderClient.UpdateOrder(c.Request.Context(), updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}

	return updatedOrder, true
}

func (h *Handler) deleteOrder(c *gin.Context, orderID uint64) bool {
	resp, err := h.orderClient.DeleteOrder(c.Request.Context(), &pbOrders.DeleteOrderRequest{
		Id: orderID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}

	if !resp.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return false
	}

	return true
}

func (h *Handler) getValidatedOrder(c *gin.Context) (*pbOrders.OrderResponse, bool) {
	userID, ok := h.getUserID(c)
	if !ok {
		return nil, false
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return nil, false
	}

	order, err := h.orderClient.GetOrder(c.Request.Context(), &pbOrders.GetOrderRequest{Id: orderID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return nil, false
	}

	if order.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return nil, false
	}

	return order, true
}
