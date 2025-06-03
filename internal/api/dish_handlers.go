package api

import (
	"net/http"

	pbDishes "github.com/anyviewww/bff-service/proto/dishes"
	"github.com/gin-gonic/gin"
)

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
