package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chiangwanyi/iot_test/db"
	"github.com/chiangwanyi/iot_test/routes"
	"github.com/chiangwanyi/iot_test/server"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始数据库
	db.Init()

	// 启动TCP服务
	server, err := server.NewTcpServer(":9309")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	go server.StartTcpServer()

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
	routes.SetupRoutes(router, db.SqliteConn)

	// 启动服务器
	port := "8080"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// 启动HTTP服务器
	go func() {
		log.Printf("服务器启动在端口 %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 中断信号处理
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Received interrupt signal, shutting down...")

	// 创建5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP服务器关闭失败: %v", err)
	}
	log.Println("HTTP服务器已优雅关闭")

	server.StopTcpServer()
	log.Println("TCP服务器已优雅关闭")
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
