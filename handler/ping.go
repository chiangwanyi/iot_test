package handler

import "github.com/gin-gonic/gin"

// PingHandler 处理 /ping 路由请求
func PingHandler(c *gin.Context) {
    c.String(200, "pong")
}