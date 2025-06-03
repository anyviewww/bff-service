package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

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
