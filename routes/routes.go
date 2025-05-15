package routes

import (
	"database/sql"

	"github.com/chiangwanyi/iot_test/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// // 创建模型实例
	// userModel := &models.UserModel{DB: db}
	// deviceModel := &models.DeviceModel{DB: db}

	// // 创建处理器实例
	// authHandler := &handlers.AuthHandler{
	// 	UserModel: userModel,
	// 	SecretKey: []byte("your-secret-key"), // 实际应用中应从配置文件或环境变量获取
	// }

	// deviceHandler := &handlers.DeviceHandler{
	// 	DeviceModel: deviceModel,
	// }

	// 公开路由组
	// public := router.Group("/api")
	// {
	// 	// 认证相关路由
	// 	public.POST("/login", authHandler.Login)
	// 	public.POST("/register", authHandler.Register)
	// }

	// // 认证后路由组
	// authenticated := router.Group("/api")
	// authenticated.Use(authHandler.AuthMiddleware())
	// {
	// 	// 设备管理路由
	// 	authenticated.POST("/devices", deviceHandler.CreateDevice)
	// 	authenticated.GET("/devices/:device_id", deviceHandler.GetDevice)
	// 	authenticated.PUT("/devices/:device_id/status", deviceHandler.UpdateDeviceStatus)
	// 	authenticated.POST("/devices/:device_id/data", deviceHandler.SaveDeviceData)
	// 	authenticated.GET("/devices/:device_id/data", deviceHandler.GetDeviceData)
	// }

	// // 管理员路由组
	// admin := router.Group("/api/admin")
	// admin.Use(authHandler.AuthMiddleware(), authHandler.AdminMiddleware())
	// {
	// 	// 管理员专用路由
	// }

	// 添加 /ping 路由
	router.GET("/ping", handler.PingHandler)
}
