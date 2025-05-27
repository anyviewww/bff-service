package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anyviewww/bff-service/internal/config"
	pbDishes "github.com/anyviewww/bff-service/proto/dishes"
	pbOrders "github.com/anyviewww/bff-service/proto/orders"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	menuClient  pbDishes.DishServiceClient
	orderClient pbOrders.OrderServiceClient
	cfg         *config.Config
}

// Basic structure for queries
type orderRequest struct {
	Items  []int64 `json:"items,omitempty"`
	Status *string `json:"status,omitempty"`
}

// Initialisation
func NewHandler(menuClient pbDishes.DishServiceClient, orderClient pbOrders.OrderServiceClient, cfg *config.Config) *Handler {
	return &Handler{
		menuClient:  menuClient,
		orderClient: orderClient,
		cfg:         cfg,
	}
}

// Menu Handlers
func (h *Handler) GetDish(c *gin.Context) {
	id, ok := h.parseIDParam(c, "id", 32)
	if !ok {
		return
	}

	dish, ok := h.getDish(c, int32(id))
	if !ok {
		return
	}

	c.JSON(http.StatusOK, h.toDishResponse(dish))
}

func (h *Handler) GetAllDishes(c *gin.Context) {
	dishes, ok := h.getAllDishes(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"dishes": dishes})
}

// Order Handlers
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

// Supporting methods
func (h *Handler) getValidatedOrder(c *gin.Context) (*pbOrders.OrderResponse, bool) {
	userID, ok := h.getUserID(c)
	if !ok {
		return nil, false
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if !ok {
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

func (h *Handler) validateUserAccess(c *gin.Context, paramUserID string) (uint64, bool) {
	tokenUserID, ok := h.getUserID(c)
	if !ok {
		return 0, false
	}

	userID, err := strconv.ParseUint(paramUserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return 0, false
	}

	if tokenUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return 0, false
	}

	return userID, true
}

// Basic operations
func (h *Handler) parseIDParam(c *gin.Context, param string, bitSize int) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(param), 10, bitSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid %s format", param)})
		return 0, false
	}
	return id, true
}

func (h *Handler) bindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func (h *Handler) getUserID(c *gin.Context) (uint64, bool) {
	claims, exists := c.Get("jwtClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return 0, false
	}

	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return 0, false
	}

	userID, ok := jwtClaims["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return 0, false
	}

	return uint64(userID), true
}

// Methods for dealing with the dishes
func (h *Handler) getDish(c *gin.Context, id int32) (*pbDishes.Dish, bool) {
	resp, err := h.menuClient.GetDishes(c.Request.Context(), &pbDishes.DishRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}

	if len(resp.Dishes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return nil, false
	}

	return resp.Dishes[0], true
}

func (h *Handler) getAllDishes(c *gin.Context) ([]gin.H, bool) {
	resp, err := h.menuClient.GetDishes(c.Request.Context(), &pbDishes.DishRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}

	dishes := make([]gin.H, 0, len(resp.Dishes))
	for _, dish := range resp.Dishes {
		dishes = append(dishes, h.toDishResponse(dish))
	}

	return dishes, true
}

// Methods for working with orders
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

// Conversion to API response
func (h *Handler) toDishResponse(dish *pbDishes.Dish) gin.H {
	return gin.H{
		"id":       dish.Id,
		"name":     dish.Name,
		"type":     gin.H{"id": dish.Type.Id, "name": dish.Type.TypeDish},
		"category": gin.H{"id": dish.Category.Id, "name": dish.Category.CategoryDish},
		"nutrition": gin.H{
			"calories":      dish.NutFact.Calories,
			"proteins":      dish.NutFact.Proteins,
			"fats":          dish.NutFact.Fats,
			"carbohydrates": dish.NutFact.Carbohydrates,
		},
		"tag":    gin.H{"id": dish.Tag.Id, "name": dish.Tag.TagDish},
		"recipe": dish.Recipe,
	}
}

func (h *Handler) toOrderResponse(order *pbOrders.OrderResponse) gin.H {
	return gin.H{
		"id":      order.Id,
		"user_id": order.UserId,
		"items":   order.Items,
		"status":  order.Status,
	}
}

func (h *Handler) getUserIDFromToken(c *gin.Context) (uint64, error) {
	claims, exists := c.Get("jwtClaims")
	if !exists {
		return 0, errors.New("token claims not found")
	}

	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims format")
	}

	userID, ok := jwtClaims["id"].(float64)
	if !ok {
		return 0, errors.New("user ID not found in token")
	}

	return uint64(userID), nil
}
