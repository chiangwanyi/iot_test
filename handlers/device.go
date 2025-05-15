package handlers

import (
	"github.com/chiangwanyi/iot_test/models"
	"github.com/gin-gonic/gin"
)

// DeviceHandler 设备处理器
type DeviceHandler struct {
	DeviceModel *models.DeviceModel
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.DeviceModel.CreateDevice(&device); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 执行成功后返回 code 为 200 的 JSON
	c.JSON(200, gin.H{"code": 200, "message": "创建设备成功！"})
}

// GetAllDevices 处理获取所有设备列表的请求
func (h *DeviceHandler) GetAllDevices(c *gin.Context) {
	devicesResp, err := h.DeviceModel.GetAllDevices()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, devicesResp)
}
