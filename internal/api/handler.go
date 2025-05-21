package api

import (
	"errors"
	"net/http"
	"strconv"

	pbDishes "github.com/anyviewww/bff-service/proto/dishes"
	pbOrders "github.com/anyviewww/bff-service/proto/orders"

	"github.com/anyviewww/bff-service/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	menuClient  pbDishes.DishServiceClient
	orderClient pbOrders.OrderServiceClient
	cfg         *config.Config
}

func NewHandler(menuClient pbDishes.DishServiceClient, orderClient pbOrders.OrderServiceClient) *Handler {
	return &Handler{
		menuClient:  menuClient,
		orderClient: orderClient,
	}
}

func (h *Handler) SetConfig(cfg *config.Config) {
	h.cfg = cfg
}

// Menu Handlers

func (h *Handler) GetDish(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dish ID format"})
		return
	}

	resp, err := h.menuClient.GetDishes(c.Request.Context(), &pbDishes.DishRequest{
		Id: int32(id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(resp.Dishes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	c.JSON(http.StatusOK, toDishResponse(resp.Dishes[0]))
}

func (h *Handler) GetAllDishes(c *gin.Context) {
	resp, err := h.menuClient.GetDishes(c.Request.Context(), &pbDishes.DishRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dishes := make([]gin.H, 0, len(resp.Dishes))
	for _, dish := range resp.Dishes {
		dishes = append(dishes, toDishResponse(dish))
	}

	c.JSON(http.StatusOK, gin.H{"dishes": dishes})
}

func toDishResponse(dish *pbDishes.Dish) gin.H {
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

// Order Handlers
func (h *Handler) CreateOrder(c *gin.Context) {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Items []int64 `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderClient.CreateOrder(c.Request.Context(), &pbOrders.CreateOrderRequest{
		UserId: userID, // Используем userID из токена
		Items:  req.Items,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toOrderResponse(order))
}
func (h *Handler) GetUserOrders(c *gin.Context) {
	// Получаем userID из токена
	tokenUserID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Получаем userID из URL параметра
	paramUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Проверяем, что пользователь запрашивает свои заказы
	if tokenUserID != paramUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own orders"})
		return
	}

}
func (h *Handler) GetOrder(c *gin.Context) {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	order, err := h.orderClient.GetOrder(c.Request.Context(), &pbOrders.GetOrderRequest{
		Id: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order))
}

func (h *Handler) UpdateOrder(c *gin.Context) {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	currentOrder, err := h.orderClient.GetOrder(c.Request.Context(), &pbOrders.GetOrderRequest{
		Id: orderID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if currentOrder.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own orders"})
		return
	}

	var req struct {
		Items  []int64 `json:"items,omitempty"`
		Status *string `json:"status,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(updatedOrder))
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	currentOrder, err := h.orderClient.GetOrder(c.Request.Context(), &pbOrders.GetOrderRequest{
		Id: orderID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if currentOrder.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own orders"})
		return
	}

	resp, err := h.orderClient.DeleteOrder(c.Request.Context(), &pbOrders.DeleteOrderRequest{
		Id: orderID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.Deleted {
		c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
	}
}

func toOrderResponse(order *pbOrders.OrderResponse) gin.H {
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

	userID, ok := jwtClaims["id"].(float64) // JWT числа всегда float64
	if !ok {
		return 0, errors.New("user ID not found in token")
	}

	return uint64(userID), nil
}
