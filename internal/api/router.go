package api

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{handler: handler}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	{
		// Menu endpoints
		menu := api.Group("/menu")
		{
			menu.GET("/dishes/:id", r.handler.GetDish)
		}

		// Order endpoints
		orders := api.Group("/orders")
		{
			orders.POST("/", r.handler.CreateOrder)
			orders.GET("/:id", r.handler.GetOrder)
			orders.PUT("/:id", r.handler.UpdateOrder)
			orders.DELETE("/:id", r.handler.DeleteOrder)
		}
	}

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
