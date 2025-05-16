package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chiangwanyi/iot_test/db"
	"github.com/chiangwanyi/iot_test/routes"
	"github.com/chiangwanyi/iot_test/tcp_mgr"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化TCP管理器
	tcpMgr := tcp_mgr.NewTcpMgr()
	// 启动TCP服务
	go func() {
		if err := tcpMgr.Start(); err != nil {
			log.Printf("启动TCP服务失败: %v", err)
		}
	}()

	// 初始数据库
	db.Init()

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建默认路由引擎
	router := gin.Default()

	// 配置中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// 注册路由
	routes.SetupRoutes(router, db.SqliteConn, tcpMgr)

	// 启动服务器
	port := "8080"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("服务器启动在端口 %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// 跨域中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
