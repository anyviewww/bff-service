package api

import (
	"bff-service/grpc/menu"
	pb "bff-service/grpc/menu"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetDishesRequest struct {
	ID int32 `json:"id"`
}

var baseHandler = BaseHandler{}

func GetDishes(c *gin.Context) {
	var req GetDishesRequest
	if !baseHandler.BindAndValidate(c, &req) {
		return
	}

	ctx := context.Background()
	resp, err := menu.GetDishes(ctx, &pb.DishRequest{
		Id: req.ID,
	})
	baseHandler.HandleGRPCError(c, err)

	dishes := make([]map[string]interface{}, len(resp.Dishes))
	for i, dish := range resp.Dishes {
		dishes[i] = map[string]interface{}{
			"id":       dish.Id,
			"name":     dish.Name,
			"type":     dish.Type.TypeDish,
			"category": dish.Category.CategoryDish,
			"nutrition_fact": map[string]interface{}{
				"calories":      dish.NutFact.Calories,
				"proteins":      dish.NutFact.Proteins,
				"fats":          dish.NutFact.Fats,
				"carbohydrates": dish.NutFact.Carbohydrates,
			},
			"tag":    dish.Tag.TagDish,
			"recipe": dish.Recipe,
		}
	}

	baseHandler.Respond(c, http.StatusOK, gin.H{"dishes": dishes})
}
