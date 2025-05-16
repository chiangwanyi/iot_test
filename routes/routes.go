package routes

import (
	"database/sql"

	"github.com/chiangwanyi/iot_test/handlers"
	"github.com/chiangwanyi/iot_test/models"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// 创建模型实例
	deviceModel := &models.DeviceModel{DB: db}
	deviceModel.CreateTables()

	deviceHandler := &handlers.DeviceHandler{
		DeviceModel: deviceModel,
	}

	router.GET("/ping", handlers.PingHandler)

	public := router.Group("/api")
	{
		// 设备管理路由
		public.POST("/devices", deviceHandler.CreateDevice)
		// 新增获取所有设备的路由
		public.GET("/devices/list", deviceHandler.GetAllDevices)
		// 新增分页查询设备的路由
		public.GET("/devices/page", deviceHandler.GetDevicesWithPage)
	}
}
