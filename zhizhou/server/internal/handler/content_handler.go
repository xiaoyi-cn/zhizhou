package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-username/zhizhou/server/internal/service"
)

type ContentHandler struct {
	ingestSvc  service.IngestService
	contentSvc service.ContentService
}

func NewContentHandler(ingestSvc service.IngestService, contentSvc service.ContentService) *ContentHandler {
	return &ContentHandler{
		ingestSvc:  ingestSvc,
		contentSvc: contentSvc,
	}
}

func (h *ContentHandler) Ingest(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	item, err := h.ingestSvc.ProcessURL(userID, req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     item.ID,
		"status": item.Status,
	})
}

func (h *ContentHandler) GetPending(c *gin.Context) {
	userID := c.GetString("userID")
	contents, err := h.contentSvc.GetPending(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get pending contents"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"contents": contents})
}

func (h *ContentHandler) Approve(c *gin.Context) {
	id := c.Param("id")
	if err := h.contentSvc.Approve(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *ContentHandler) Skip(c *gin.Context) {
	id := c.Param("id")
	if err := h.contentSvc.Skip(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to skip"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "skipped"})
}

func (h *ContentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title    string   `json:"title"`
		Summary  string   `json:"summary"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.contentSvc.Update(id, req.Title, req.Summary, req.Category, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *ContentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	content, err := h.contentSvc.GetByID(id)
	if err != nil || content == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": content})
}

func (h *ContentHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	tags := c.QueryArray("tags")
	category := c.Query("category")

	contents, total, err := h.contentSvc.List(userID, tags, category, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list contents"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (h *ContentHandler) Search(c *gin.Context) {
	userID := c.GetString("userID")
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	tags := c.QueryArray("tags")

	contents, total, err := h.contentSvc.Search(userID, query, tags, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"query":    query,
	})
}