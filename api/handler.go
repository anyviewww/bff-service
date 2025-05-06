package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BaseHandler represents the base handler for all API methods.
type BaseHandler struct{}

// Respond sends a JSON response to the client.
func (h *BaseHandler) Respond(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// BindAndValidate performs binding of JSON data to a structure and validation.
func (h *BaseHandler) BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.Respond(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return false
	}
	return true
}

// HandleGRPCError handles errors that occur when calling gRPC services.
func (h *BaseHandler) HandleGRPCError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
}
