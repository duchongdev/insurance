package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/repository"
	"github.com/huaan/insurance-bridge/internal/service"
)

// AdminHandler 管理后台 REST API：登录、渠道 CRUD、银行信息查询与刷新、华安配置。
type AdminHandler struct {
	admin     *service.AdminService
	proxy     *service.ProxyService
	huaAnConf *service.HuaAnConfigService
}

// NewAdminHandler 构造管理后台处理器。
func NewAdminHandler(admin *service.AdminService, proxy *service.ProxyService, huaAnConf *service.HuaAnConfigService) *AdminHandler {
	return &AdminHandler{admin: admin, proxy: proxy, huaAnConf: huaAnConf}
}

// Register 在 /admin/api 路由组下注册所有管理接口；除 login 外均需 Bearer JWT。
func (h *AdminHandler) Register(r *gin.RouterGroup) {
	r.POST("/login", h.login)
	auth := r.Group("")
	auth.Use(h.authMiddleware)
	auth.GET("/channels", h.listChannels)
	auth.POST("/channels", h.createChannel)
	auth.PUT("/channels/:id", h.updateChannel)
	auth.DELETE("/channels/:id", h.deleteChannel)
	auth.GET("/banks", h.listBanks)
	auth.POST("/bank-list/refresh", h.refreshBankList)
	auth.GET("/bank-list", h.getBankList)
	auth.GET("/huaan-configs", h.listHuaAnConfigs)
	auth.POST("/huaan-configs", h.createHuaAnConfig)
	auth.PUT("/huaan-configs/:id", h.updateHuaAnConfig)
	auth.DELETE("/huaan-configs/:id", h.deleteHuaAnConfig)
	auth.POST("/huaan-configs/:id/test", h.testHuaAnConfigByID)
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

// listChannels 分页列出渠道配置。
func (h *AdminHandler) listChannels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.admin.ListChannels(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	for i := range list {
		h.enrichChannel(&list[i])
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// enrichChannel 填充关联华安配置的密钥与环境类型。
func (h *AdminHandler) enrichChannel(ch *model.Channel) {
	if ch == nil || ch.HuaAnSettingID == 0 {
		return
	}
	conf, err := h.huaAnConf.GetByID(ch.HuaAnSettingID)
	if err != nil {
		return
	}
	ch.HuaAnKey = conf.ChannelSecret
	ch.EnvType = conf.EnvType
}

// syncChannelHuaAnKey 若渠道关联华安配置，则同步展示用 huaAnKey。
func (h *AdminHandler) syncChannelHuaAnKey(ch *model.Channel) {
	h.enrichChannel(ch)
}

// createChannel 创建渠道；渠道编码与华安密钥来自所选华安配置。
func (h *AdminHandler) createChannel(c *gin.Context) {
	var req struct {
		HuaAnSettingID uint64 `json:"huaanSettingId"`
		ChannelName    string `json:"channelName"`
		ChannelKey     string `json:"channelKey"`
		CallbackURL    string `json:"callbackUrl"`
		Status         int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if req.HuaAnSettingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请选择华安配置中的渠道编码"})
		return
	}
	code, secret, err := h.huaAnConf.ResolveChannelCredentials(req.HuaAnSettingID)
	if err != nil {
		if errors.Is(err, service.ErrHuaAnConfigNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "所选华安配置不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ch := &model.Channel{
		HuaAnSettingID: req.HuaAnSettingID,
		ChannelCode:    code,
		HuaAnKey:       secret,
		ChannelName:    req.ChannelName,
		ChannelKey:     req.ChannelKey,
		CallbackURL:    req.CallbackURL,
		Status:         req.Status,
	}
	if ch.Status == 0 {
		ch.Status = 1
	}
	if err := h.admin.CreateChannel(ch); err != nil {
		if errors.Is(err, service.ErrInvalidChannelConfig) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	h.enrichChannel(ch)
	c.JSON(http.StatusOK, ch)
}

// updateChannel 更新渠道名称、回调地址、状态与渠道密钥。
func (h *AdminHandler) updateChannel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		ChannelName string `json:"channelName"`
		ChannelKey  string `json:"channelKey"`
		CallbackURL string `json:"callbackUrl"`
		Status      int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	ch, err := h.admin.UpdateChannel(id, &model.Channel{
		ChannelName: req.ChannelName,
		ChannelKey:  req.ChannelKey,
		CallbackURL: req.CallbackURL,
		Status:      req.Status,
	})
	if err != nil {
		if errors.Is(err, repository.ErrChannelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "channel not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidChannelConfig) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	h.syncChannelHuaAnKey(ch)
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

// refreshBankList 按所选华安配置请求上游 getBankList。
func (h *AdminHandler) refreshBankList(c *gin.Context) {
	var req struct {
		HuaAnSettingID uint64 `json:"huaanSettingId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if req.HuaAnSettingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请选择华安配置"})
		return
	}
	out, err := h.proxy.RefreshBankList(c.Request.Context(), req.HuaAnSettingID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "所选华安配置不存在"})
			return
		}
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

// listBanks 分页查询 bank_info_t 银行信息。
func (h *AdminHandler) listBanks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	list, total, err := h.admin.ListBanks(int8(status), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// listHuaAnConfigs 返回全部华安上游配置。
func (h *AdminHandler) listHuaAnConfigs(c *gin.Context) {
	list, err := h.huaAnConf.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list})
}

// createHuaAnConfig 新建华安配置。
func (h *AdminHandler) createHuaAnConfig(c *gin.Context) {
	var req service.HuaAnConfigDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	conf, err := h.huaAnConf.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidHuaAnConfig) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请填写环境类型、接口地址、渠道编码与密钥"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conf)
}

// updateHuaAnConfig 更新华安配置。
func (h *AdminHandler) updateHuaAnConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req service.HuaAnConfigDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	conf, err := h.huaAnConf.Update(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrHuaAnConfigNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "config not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidHuaAnConfig) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请填写环境类型、接口地址与渠道编码"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conf)
}

// deleteHuaAnConfig 删除华安配置。
func (h *AdminHandler) deleteHuaAnConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.huaAnConf.Delete(id); err != nil {
		if errors.Is(err, service.ErrHuaAnConfigNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "config not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// testHuaAnConfigByID 对指定配置的接口地址做 HTTP 连通性探测。
func (h *AdminHandler) testHuaAnConfigByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := h.huaAnConf.TestConnectivityByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrHuaAnConfigNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "config not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidHuaAnConfig) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请填写有效的接口地址"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
