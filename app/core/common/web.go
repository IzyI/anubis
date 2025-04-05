package common

import (
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

func GetClientIP(ctx *gin.Context) string {
	ip := ctx.GetHeader("X-Real-IP")
	if ip == "" {
		ip = ctx.GetHeader("X-Forwarded-For")
	}

	if ip != "" {
		ips := strings.Split(ip, ",")
		ip = strings.TrimSpace(ips[0])
	}

	if ip == "" {
		ip, _, _ = net.SplitHostPort(ctx.Request.RemoteAddr)
	}

	return ip
}

func GetDeviceType(ctx *gin.Context) string {
	deviceType := ctx.GetHeader("User-Agent")
	if deviceType == "" {
		// Если заголовок X-Device не присутствует, используем User-Agent
		deviceType = ctx.GetHeader("X-Device")
	} else {
		deviceType = "_"
	}
	return deviceType
}
