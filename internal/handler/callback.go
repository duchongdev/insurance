package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huaan/insurance-bridge/internal/callback"
	"github.com/huaan/insurance-bridge/internal/service"
)

// CallbackHandler 华安调用本服务的回调入口（投保结果通知等）。
type CallbackHandler struct {
	proxy *service.ProxyService
}

// NewCallbackHandler 构造回调处理器。
func NewCallbackHandler(proxy *service.ProxyService) *CallbackHandler {
	return &CallbackHandler{proxy: proxy}
}

// Register 注册华安 → 本服务回调路由（独立于渠道 upChannelApi）。
func (h *CallbackHandler) Register(r gin.IRouter) {
	r.POST(callback.InsureNotifyPath, h.insureNotify)
}

func (h *CallbackHandler) insureNotify(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "read body failed"})
		return
	}
	out, status, _ := h.proxy.HandleInsureNotify(c.Request.Context(), raw)
	c.Data(status, "application/json; charset=utf-8", out)
}
