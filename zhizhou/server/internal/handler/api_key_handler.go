package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/zhizhou/server/internal/service"
)

type APIKeyHandler struct {
	service service.APIKeyService
}

func NewAPIKeyHandler(s service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{service: s}
}

func (h *APIKeyHandler) Save(c *gin.Context) {
	var req struct {
		Provider string `json:"provider" binding:"required"`
		Key      string `json:"key" binding:"required"`
		BaseURL  string `json:"base_url"`
		Model    string `json:"model" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	err := h.service.SaveAPIKey(userID, req.Provider, req.Key, req.BaseURL, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *APIKeyHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	keys, err := h.service.ListAPIKeys(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list keys"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func (h *APIKeyHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Provider string `json:"provider" binding:"required"`
		Key      string `json:"key" binding:"required"`
		BaseURL  string `json:"base_url"`
		Model    string `json:"model" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateAPIKey(id, req.Provider, req.Key, req.BaseURL, req.Model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *APIKeyHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteAPIKey(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *APIKeyHandler) Test(c *gin.Context) {
	var req struct {
		Provider string `json:"provider" binding:"required"`
		Key      string `json:"key" binding:"required"`
		BaseURL  string `json:"base_url"`
		Model    string `json:"model" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.TestAPIKey(req.Provider, req.Key, req.BaseURL, req.Model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "connection test failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "connection successful"})
}