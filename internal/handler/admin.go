package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/service"
)

// AdminHandler 管理后台 REST API：登录、渠道 CRUD、审计日志与业务数据分页查询。
type AdminHandler struct {
	admin *service.AdminService
	proxy *service.ProxyService
}

// NewAdminHandler 构造管理后台处理器。
func NewAdminHandler(admin *service.AdminService, proxy *service.ProxyService) *AdminHandler {
	return &AdminHandler{admin: admin, proxy: proxy}
}

// Register 在 /admin/api 路由组下注册所有管理接口；除 login 外均需 Bearer JWT。
func (h *AdminHandler) Register(r *gin.RouterGroup) {
	r.POST("/login", h.login)
	auth := r.Group("")
	auth.Use(h.authMiddleware)
	auth.GET("/stats", h.stats)
	auth.GET("/channels", h.listChannels)
	auth.POST("/channels", h.createChannel)
	auth.PUT("/channels/:id", h.updateChannel)
	auth.DELETE("/channels/:id", h.deleteChannel)
	auth.GET("/logs", h.listLogs)
	auth.GET("/policies", h.listPolicies)
	auth.GET("/users", h.listUsers)
	auth.GET("/signs", h.listSigns)
	auth.POST("/bank-list/refresh", h.refreshBankList)
	auth.GET("/bank-list", h.getBankList)
}

// authMiddleware 从 Authorization: Bearer <token> 解析 JWT，失败返回 401。
func (h *AdminHandler) authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	if _, err := h.admin.ParseToken(token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}
	c.Next()
}

// login 用户名密码登录，成功返回 JWT（24h 有效）。
func (h *AdminHandler) login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	token, err := h.admin.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "login failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// stats 按渠道（可选 channelCode 查询参数）聚合保单/用户/签约/日志数量。
func (h *AdminHandler) stats(c *gin.Context) {
	stats, err := h.admin.Stats(c.Query("channelCode"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// listChannels 分页列出渠道配置。
func (h *AdminHandler) listChannels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.admin.ListChannels(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// createChannel 创建渠道；未传 channelKey 时自动生成。
func (h *AdminHandler) createChannel(c *gin.Context) {
	var ch model.Channel
	if err := c.ShouldBindJSON(&ch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if err := h.admin.CreateChannel(&ch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ch)
}

// updateChannel 按路径 id 更新渠道（请求体需含可更新字段）。
func (h *AdminHandler) updateChannel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var ch model.Channel
	if err := c.ShouldBindJSON(&ch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	ch.ID = id
	if err := h.admin.UpdateChannel(&ch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ch)
}

// deleteChannel 按 id 删除渠道。
func (h *AdminHandler) deleteChannel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.admin.DeleteChannel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// listLogs 分页查询接口审计日志，支持 channelCode、apiPath 过滤。
func (h *AdminHandler) listLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.admin.ListLogs(c.Query("channelCode"), c.Query("apiPath"), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// listPolicies 分页查询抽取的保单记录。
func (h *AdminHandler) listPolicies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.admin.ListPolicies(c.Query("channelCode"), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// listUsers 分页查询用户记录（三要素为库内密文）。
func (h *AdminHandler) listUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.admin.ListUsers(c.Query("channelCode"), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// refreshBankList 使用指定渠道码与华安密钥请求上游 getBankList。
func (h *AdminHandler) refreshBankList(c *gin.Context) {
	var req struct {
		ChannelCode string `json:"channelCode"`
		HuaAnKey    string `json:"huaAnKey"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if req.ChannelCode == "" || req.HuaAnKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "channelCode and huaAnKey required"})
		return
	}
	out, err := h.proxy.RefreshBankList(c.Request.Context(), req.ChannelCode, req.HuaAnKey)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", out)
}

// getBankList 读取 Redis 中已缓存的银行列表（华安原始 JSON）。
func (h *AdminHandler) getBankList(c *gin.Context) {
	channelCode := c.Query("channelCode")
	if channelCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "channelCode required"})
		return
	}
	data, err := h.proxy.GetCachedBankList(c.Request.Context(), channelCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "bank list not cached"})
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

// listSigns 分页查询签约记录。
func (h *AdminHandler) listSigns(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.admin.ListSigns(c.Query("channelCode"), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}
