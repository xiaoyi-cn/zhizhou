package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/zhizhou/server/internal/middleware"
	"github.com/your-username/zhizhou/server/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SendCode(req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "code sent"})
}

func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.service.VerifyCode(req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "verification failed"})
		return
	}

	token, err := middleware.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user_id": userID,
	})
}