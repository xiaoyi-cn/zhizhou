package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// 需要专业版的功能列表
var proFeatures = map[string]bool{
	"mcp_server": true,
	"open_api":   true,
	"auto_rules": true,
}

// TierRequired 版本控制中间件
// feature: 功能名称，如 "mcp_server", "open_api"
func TierRequired(feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 本地部署跳过所有权限检查
		if os.Getenv("DEPLOY_MODE") == "local" {
			c.Next()
			return
		}

		// 检查是否是需要专业版的功能
		if proFeatures[feature] {
			tier := c.GetString("tier")
			if tier != "pro" {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "FeatureRequiresPro",
					"message": "此功能需要专业版订阅",
					"feature": feature,
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// StorageQuotaRequired 存储配额检查中间件
func StorageQuotaRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("DEPLOY_MODE") == "local" {
			c.Next()
			return
		}

		// 在 ingest handler 中单独检查
		c.Next()
	}
}