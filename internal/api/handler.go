import (
	"net/http"
	"strconv"

	pbDishes "github.com/anyviewww/bff-service/proto/dishes"
	pbOrders "github.com/anyviewww/bff-service/proto/orders"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	menuClient  pbDishes.DishServiceClient
	orderClient pbOrders.OrderServiceClient
}

func NewHandler(menuClient pbDishes.DishServiceClient, orderClient pbOrders.OrderServiceClient) *Handler {
	return &Handler{
		menuClient:  menuClient,
		orderClient: orderClient,
	}
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
	var req struct {
		UserID uint64  `json:"user_id" binding:"required"`
		Items  []int64 `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderClient.CreateOrder(c.Request.Context(), &pbOrders.CreateOrderRequest{
		UserId: req.UserID,
		Items:  req.Items,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toOrderResponse(order))
}

func (h *Handler) GetOrder(c *gin.Context) {
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	var req struct {
		UserID *uint64 `json:"user_id,omitempty"`
		Items  []int64 `json:"items,omitempty"`
		Status *string `json:"status,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := &pbOrders.UpdateOrderRequest{Id: id}
	if req.UserID != nil {
		updateReq.UserId = *req.UserID
	}
	if req.Items != nil {
		updateReq.Items = req.Items
	}
	if req.Status != nil {
		updateReq.Status = *req.Status
	}

	order, err := h.orderClient.UpdateOrder(c.Request.Context(), updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order))
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	resp, err := h.orderClient.DeleteOrder(c.Request.Context(), &pbOrders.DeleteOrderRequest{
		Id: id,
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

func (h *Handler) GetUserOrders(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// В текущей реализации proto нет метода для получения заказов пользователя
	// Это пример того, как можно реализовать, если добавить метод в OrderService
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func toOrderResponse(order *pbOrders.OrderResponse) gin.H {
	return gin.H{
		"id":      order.Id,
		"user_id": order.UserId,
		"items":   order.Items,
		"status":  order.Status,
	}
}
