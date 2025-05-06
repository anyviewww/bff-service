package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

var baseHandler = BaseHandler{}

func Login(c *gin.Context) {
	var req LoginRequest
	if !baseHandler.BindAndValidate(c, &req) {
		return
	}
	baseHandler.Respond(c, http.StatusOK, gin.H{"message": "Login successful", "email": req.Email})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if !baseHandler.BindAndValidate(c, &req) {
		return
	}
	baseHandler.Respond(c, http.StatusOK, gin.H{"message": "User registered successfully", "email": req.Email})
}

func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if !baseHandler.BindAndValidate(c, &req) {
		return
	}
	baseHandler.Respond(c, http.StatusOK, gin.H{"message": "Password reset email sent", "email": req.Email})
}
