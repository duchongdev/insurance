module github.com/huaan/insurance-bridge

// 本地开发需 Go 1.22+（依赖标准库 log/slog、slices 及 gin/viper 等间接依赖）
go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/viper v1.19.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.28.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.11
)
