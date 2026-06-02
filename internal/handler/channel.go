// Package handler 注册 HTTP 路由：渠道代理 API、管理后台 API、健康检查与静态资源。
package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huaan/insurance-bridge/internal/huaan"
	"github.com/huaan/insurance-bridge/internal/service"
)

// ChannelHandler 将各渠道 API 统一委托给 ProxyService.Forward 处理。
type ChannelHandler struct {
	proxy *service.ProxyService
}

// NewChannelHandler 构造渠道 API 处理器。
func NewChannelHandler(proxy *service.ProxyService) *ChannelHandler {
	return &ChannelHandler{proxy: proxy}
}

// Register 为 huaan.APIPaths 中每个路径注册 POST 处理器，挂载在配置的 HuaAn.APIPath 路由组下。
func (h *ChannelHandler) Register(r *gin.RouterGroup) {
	for _, path := range huaan.APIPaths {
		p := path
		r.POST(p, func(c *gin.Context) {
			h.handle(c, p)
		})
	}
}

// handle 读取原始 JSON body 并转发至代理层；HTTP status 由 ProxyService 固定为 200。
func (h *ChannelHandler) handle(c *gin.Context, apiPath string) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "read body failed", "data": nil})
		return
	}
	out, status, _ := h.proxy.Forward(c.Request.Context(), apiPath, raw)
	c.Data(status, "application/json; charset=utf-8", out)
}
