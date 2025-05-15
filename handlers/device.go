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
}
