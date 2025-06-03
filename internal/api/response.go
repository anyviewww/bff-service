package api

import (
	pbDishes "github.com/anyviewww/bff-service/proto/dishes"
	pbOrders "github.com/anyviewww/bff-service/proto/orders"
	"github.com/gin-gonic/gin"
)

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
