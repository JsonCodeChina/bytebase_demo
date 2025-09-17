package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shenbo/sql-review-learning-demo/pkg/advisor"
	"github.com/shenbo/sql-review-learning-demo/pkg/api"
	"github.com/shenbo/sql-review-learning-demo/pkg/config"
	"github.com/shenbo/sql-review-learning-demo/pkg/database"
	"github.com/shenbo/sql-review-learning-demo/pkg/rules/mysql"
)

func main() {
	// 加载配置
	loader := config.NewLoader("config")
	cfg, err := loader.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建数据库管理器
	dbManager := database.NewDatabaseManagerWithConfig(cfg.Database)
	defer dbManager.Close()

	// 创建审查器并注册规则
	sqlAdvisor := advisor.NewDefaultAdvisor()
	sqlAdvisor.RegisterRule(mysql.NewTableRequirePKRule())

	// 创建HTTP服务器
	server := api.NewServer(dbManager, sqlAdvisor)

	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", strings.Join(cfg.Server.CORS.AllowOrigins, ","))
		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.Server.CORS.AllowMethods, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.Server.CORS.AllowHeaders, ","))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API路由
	api := r.Group("/api")
	{
		// 连接管理
		api.POST("/connections/test", server.TestConnection)
		api.POST("/connections", server.SaveConnection)
		api.GET("/connections", server.ListConnections)

		// Schema相关
		api.GET("/schema/:connection_id", server.GetSchema)

		// SQL审查
		api.POST("/sql/review", server.ReviewSQL)

		// 规则管理
		api.GET("/rules", server.ListRules)
	}

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "sql-review-learning-demo",
		})
	})

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "SQL Review Learning Demo API Server",
			"version": "1.0.0",
			"endpoints": []string{
				"/api/connections/test",
				"/api/connections",
				"/api/schema/:connection_id",
				"/api/sql/review",
				"/api/rules",
				"/health",
			},
		})
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting SQL Review Learning Demo API Server on %s", addr)
	log.Printf("Environment: %s", loader.GetEnv())
	log.Printf("Gin Mode: %s", cfg.Server.Mode)
	log.Println("API Documentation:")
	log.Println("  POST /api/connections/test  - 测试数据库连接")
	log.Println("  POST /api/connections       - 保存数据库连接")
	log.Println("  GET  /api/connections       - 列出所有连接")
	log.Println("  GET  /api/schema/:id        - 获取数据库schema")
	log.Println("  POST /api/sql/review        - 审查SQL语句")
	log.Println("  GET  /api/rules             - 列出所有规则")

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
