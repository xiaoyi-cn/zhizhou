package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/your-username/zhizhou/server/internal/handler"
	"github.com/your-username/zhizhou/server/internal/middleware"
	"github.com/your-username/zhizhou/server/internal/repository"
	"github.com/your-username/zhizhou/server/internal/service"
)

func main() {
	// 数据库连接
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY environment variable is required")
	}

	// 依赖注入 - Repository 层
	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	contentRepo := repository.NewContentRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// 依赖注入 - Service 层
	apiKeySvc := service.NewAPIKeyService(apiKeyRepo, encryptionKey)
	authSvc := service.NewAuthService(userRepo)
	ingestSvc := service.NewIngestService(contentRepo, apiKeySvc)
	contentSvc := service.NewContentService(contentRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	subscriptionSvc := service.NewSubscriptionService(userRepo)

	// 依赖注入 - Handler 层
	authHandler := handler.NewAuthHandler(authSvc)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeySvc)
	contentHandler := handler.NewContentHandler(ingestSvc, contentSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionSvc)

	// 路由
	r := gin.Default()

	// 健康检查
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 认证路由（无需 JWT）
	auth := r.Group("/api/auth")
	{
		auth.POST("/send-code", authHandler.SendCode)
		auth.POST("/verify-code", authHandler.VerifyCode)
	}

	// 需要认证的路由
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		// API Key 管理
		api.POST("/keys", apiKeyHandler.Save)
		api.GET("/keys", apiKeyHandler.List)
		api.PUT("/keys/:id", apiKeyHandler.Update)
		api.DELETE("/keys/:id", apiKeyHandler.Delete)
		api.POST("/keys/:id/test", apiKeyHandler.Test)

		// 内容采集与消化
		api.POST("/contents/ingest", contentHandler.Ingest)
		api.GET("/contents/pending", contentHandler.GetPending)
		api.POST("/contents/:id/approve", contentHandler.Approve)
		api.POST("/contents/:id/skip", contentHandler.Skip)
		api.PUT("/contents/:id", contentHandler.Update)
		api.GET("/contents/:id", contentHandler.GetByID)
		api.GET("/contents", contentHandler.List)

		// 搜索
		api.GET("/search", contentHandler.Search)

		// 分类管理
		api.POST("/categories", categoryHandler.Create)
		api.GET("/categories", categoryHandler.List)
		api.DELETE("/categories/:id", categoryHandler.Delete)

		// 订阅与版本
		api.GET("/subscription", subscriptionHandler.GetSubscription)
		api.POST("/subscription/checkout", subscriptionHandler.Checkout)
		api.POST("/subscription/webhook", subscriptionHandler.Webhook)
		api.DELETE("/subscription", subscriptionHandler.Cancel)
		api.GET("/subscription/invoice", subscriptionHandler.GetInvoice)
		api.GET("/features", subscriptionHandler.GetFeatures)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("知舟 server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}