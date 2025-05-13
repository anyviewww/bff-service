package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Structures for dishes
type DishRequest struct {
	Name        string       `json:"name"`
	Type        DishType     `json:"type"`
	Category    DishCategory `json:"category"`
	Price       float64      `json:"price"`
	Weight      float64      `json:"weight"`
	Description string       `json:"description"`
	Nutrition   Nutrition    `json:"nutrition"`
	Tag         DishTag      `json:"tag"`
	Recipe      string       `json:"recipe"`
}

type DishType struct {
	ID int `json:"id"`
}

type DishCategory struct {
	ID int `json:"id"`
}

type Nutrition struct {
	Calories      float64 `json:"calories"`
	Proteins      float64 `json:"proteins"`
	Fats          float64 `json:"fats"`
	Carbohydrates float64 `json:"carbohydrates"`
}

type DishTag struct {
	ID int `json:"id"`
}

type MenuItem struct {
	DishID int `json:"dish_id"`
}

// Structures for orders
type Order struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Items  []int  `json:"items"`
	Status string `json:"status"`
}

type CreateOrderRequest struct {
	UserID int   `json:"user_id"`
	Items  []int `json:"items"`
}

type UpdateOrderRequest struct {
	Status *string `json:"status,omitempty"`
	Items  []int   `json:"items,omitempty"`
}

// Methods for dishes
func (s *Server) createDish(c *gin.Context) {
	var req DishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "created",
		"dish_id": 1,
	})
}

func (s *Server) createMenu(c *gin.Context) {
	var menuItems []MenuItem
	if err := c.ShouldBindJSON(&menuItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "created",
		"menu_id": 1,
		"dishes":  menuItems,
	})
}

func (s *Server) getMenu(c *gin.Context) {
	menu := []gin.H{
		{
			"dish_id": 1,
			"name":    "Шашлык",
			"price":   350.0,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"menu": menu,
	})
}

// Methods for orders
func (s *Server) getOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	orderID, _ := strconv.Atoi(id)
	c.JSON(http.StatusOK, Order{
		ID:     orderID,
		UserID: 13,
		Items:  []int{13, 15, 7},
		Status: "processing",
	})
}

func (s *Server) createOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Order{
		ID:     3,
		UserID: req.UserID,
		Items:  req.Items,
		Status: "created",
	})
}

func (s *Server) updateOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, _ := strconv.Atoi(id)
	status := "cooked"
	if req.Status != nil {
		status = *req.Status
	}

	items := []int{13, 15, 7}
	if req.Items != nil {
		items = req.Items
	}

	c.JSON(http.StatusOK, Order{
		ID:     orderID,
		UserID: 13,
		Items:  items,
		Status: status,
	})
}

func (s *Server) deleteOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("order %s deleted", id),
	})
}
