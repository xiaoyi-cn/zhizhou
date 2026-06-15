package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/your-username/zhizhou/server/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(s service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: s}
}

func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	userID := c.GetString("userID")
	info, err := h.service.GetSubscriptionInfo(userID, os.Getenv("DEPLOY_MODE"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription info"})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *SubscriptionHandler) GetFeatures(c *gin.Context) {
	userID := c.GetString("userID")
	info, err := h.service.GetSubscriptionInfo(userID, os.Getenv("DEPLOY_MODE"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get features"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tier":     info.Tier,
		"features": info.Features,
	})
}

func (h *SubscriptionHandler) Checkout(c *gin.Context) {
	var req struct {
		Plan string `json:"plan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: 接入微信支付/支付宝
	c.JSON(http.StatusOK, gin.H{"message": "checkout not yet implemented"})
}

func (h *SubscriptionHandler) Webhook(c *gin.Context) {
	// TODO: 处理支付回调
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	// TODO: 取消自动续费
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

func (h *SubscriptionHandler) GetInvoice(c *gin.Context) {
	// TODO: 获取发票
	c.JSON(http.StatusOK, gin.H{"message": "invoice not yet implemented"})
}