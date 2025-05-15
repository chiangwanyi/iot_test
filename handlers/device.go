package handlers

import (
	"fmt"

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
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	if err := h.DeviceModel.CreateDevice(&device); err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 执行成功后返回 code 为 200 的 JSON
	c.JSON(200, gin.H{"code": 200, "msg": "创建设备成功！"})
}

// GetAllDevices 处理获取所有设备列表的请求
func (h *DeviceHandler) GetAllDevices(c *gin.Context) {
	devicesResp, err := h.DeviceModel.GetAllDevices()
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "获取设备列表成功！", "data": devicesResp})
}

// GetDevicesWithPage 处理分页查询设备列表的请求
func (h *DeviceHandler) GetDevicesWithPage(c *gin.Context) {
	page := 1
	pageSize := 10

	// 从查询参数中获取页码和每页数量
	if c.Query("pageNum") != "" {
		fmt.Sscanf(c.Query("pageNum"), "%d", &page)
	} else {
		c.JSON(200, gin.H{"code": 400, "msg": "pageNum error"})
		return
	}
	if c.Query("pageSize") != "" {
		fmt.Sscanf(c.Query("pageSize"), "%d", &pageSize)
	} else {
		c.JSON(200, gin.H{"code": 400, "msg": "pageSize error"})
		return
	}

	devicesResp, totalCount, err := h.DeviceModel.GetDevicesWithPage(page, pageSize)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "获取分页设备列表成功！", "data": devicesResp, "total": totalCount})
}
