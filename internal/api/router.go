package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	s := &Server{
		router: gin.Default(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) setupRoutes() {
	// Меню
	menu := s.router.Group("/")
	{
		menu.POST("/CreateDish", s.createDish)
		menu.POST("/CreateMenu", s.createMenu)
		menu.GET("/GetMenu", s.getMenu)
	}

	// Заказы
	orders := s.router.Group("/orders")
	{
		orders.GET("/:id", s.getOrder)
		orders.POST("/", s.createOrder)
		orders.PATCH("/:id", s.updateOrder)
		orders.DELETE("/:id", s.deleteOrder)
	}
}
