// 华安保险渠道对接服务入口。
//
// 启动流程：加载配置 → 初始化日志/加密/数据库 → Seed 默认管理员 → 组装代理与管理服务 →
// 注册路由（渠道 API、管理 API、健康检查）→ 启动日志清理定时任务 → 优雅关闭。
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huaan/insurance-bridge/internal/bootstrap"
	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/handler"
	"github.com/huaan/insurance-bridge/internal/job"
	"github.com/huaan/insurance-bridge/internal/middleware"
	"github.com/huaan/insurance-bridge/internal/pkg/cipher"
	"github.com/huaan/insurance-bridge/internal/pkg/logger"
	"github.com/huaan/insurance-bridge/internal/pkg/pii"
	redisclient "github.com/huaan/insurance-bridge/internal/pkg/redis"
	"github.com/huaan/insurance-bridge/internal/repository"
	"github.com/huaan/insurance-bridge/internal/service"
	"go.uber.org/zap"
)

func main() {
	// --- 配置与基础设施 ---
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config/config.yaml"
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}

	log, err := logger.New(logger.Options{
		Level:          cfg.Log.Level,
		FilePath:       cfg.Log.FilePath,
		RetentionDays:  cfg.Log.RetentionDays,
		ArchiveEnabled: cfg.Log.ArchiveEnabled,
		MaxSizeMB:      cfg.Log.MaxSizeMB,
	})
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	crypter, err := cipher.New(cfg.Security.DataEncryptionKey)
	if err != nil {
		log.Fatal("invalid encryption key", zap.Error(err))
	}

	db, err := repository.NewDB(cfg.Database.DSN)
	if err != nil {
		log.Fatal("db connect failed", zap.Error(err))
	}

	rdb, err := redisclient.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatal("redis connect failed", zap.Error(err))
	}
	defer rdb.Close()

	// --- 仓储与首次启动 Seed ---
	channelRepo := repository.NewChannelRepo(db)
	adminRepo := repository.NewAdminRepo(db)
	bankRepo := repository.NewBankRepo(db)
	huaAnSettingRepo := repository.NewHuaAnSettingRepo(db)

	if err := bootstrap.Seed(adminRepo, log); err != nil {
		log.Fatal("seed failed", zap.Error(err))
	}
	if err := bootstrap.SyncHuaAnSetting(cfg, huaAnSettingRepo, log); err != nil {
		log.Fatal("huaan settings sync failed", zap.Error(err))
	}

	// --- 业务服务 ---
	piiTransformer := pii.NewTransformer(crypter)
	bankListCache := redisclient.NewBankListCache(rdb, cfg.Redis.BankListTTL)
	proxySvc := service.NewProxyService(cfg, log, channelRepo, bankRepo, huaAnSettingRepo, piiTransformer, bankListCache, nil)
	adminSvc := service.NewAdminService(adminRepo, channelRepo, bankRepo, cfg.Security.JWTSecret)
	huaAnConfSvc := service.NewHuaAnConfigService(cfg, huaAnSettingRepo)

	// --- HTTP 路由 ---
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Trace(), middleware.CORS(cfg.Admin.CorsOrigins))

	health := handler.NewHealthHandler(db, rdb)
	r.GET("/health/live", health.Live)
	r.GET("/health/ready", health.Ready)

	channelHandler := handler.NewChannelHandler(proxySvc)
	api := r.Group(cfg.HuaAn.APIPath)
	channelHandler.Register(api)

	callbackHandler := handler.NewCallbackHandler(proxySvc)
	callbackHandler.Register(r)

	adminHandler := handler.NewAdminHandler(adminSvc, proxySvc, huaAnConfSvc)
	adminHandler.Register(r.Group("/admin/api"))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":  "insurance-bridge",
			"admin":    "/admin/",
			"health":   "/health/ready",
			"callback": "/huaan/callback/insureNotify",
		})
	})
	r.StaticFile("/openapi.yaml", "api/openapi.yaml")

	// --- 后台任务：清理超过保留期的历史日志文件 ---
	cleanup := job.NewCleanupJob(log, cfg.Log.RetentionDays, cfg.Log.FilePath)
	cleanup.Start()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Info("server started", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen failed", zap.Error(err))
		}
	}()

	// --- 优雅关闭 ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info("server stopped")
}
