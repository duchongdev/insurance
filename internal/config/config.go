// Package config 负责加载应用配置：默认读取 config.yaml，环境变量 BRIDGE_* 可覆盖（如 BRIDGE_DATABASE_DSN）。
package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用全局配置根结构。
type Config struct {
	Server   ServerConfig   // HTTP 服务与超时
	HuaAn    HuaAnConfig    // 华安上游地址
	Database DatabaseConfig // MySQL 连接
	Redis    RedisConfig    // Redis 连接
	Security SecurityConfig // 加密与 JWT
	Log      LogConfig      // 日志与清理策略
	Admin    AdminConfig    // 默认管理员（仅首次 seed 使用）
}

// ServerConfig Gin/HTTP 服务监听与超时参数。
type ServerConfig struct {
	Port            int           // 监听端口，默认 8080
	Mode            string        // gin 模式：debug / release
	ReadTimeout     time.Duration // 读请求超时
	WriteTimeout    time.Duration // 写响应超时
	UpstreamTimeout time.Duration // 调用华安 HTTP 客户端超时
}

// HuaAnConfig 华安保险上游 API 根路径配置。
type HuaAnConfig struct {
	BaseURL string // 上游域名，如 https://xx.xxx.api.york.xin
	APIPath string // 渠道 API 前缀，默认 /upChannelApi
}

// DatabaseConfig 数据库连接。
type DatabaseConfig struct {
	DSN string // MySQL DSN
}

// RedisConfig Redis 连接参数。
type RedisConfig struct {
	Addr        string        // 地址，Docker 内默认 redis:6379
	Password    string        // 认证密码
	DB          int           // 库编号，默认 0
	BankListTTL time.Duration // 银行列表缓存 TTL，0 表示不过期
}

// SecurityConfig 安全相关密钥。
type SecurityConfig struct {
	DataEncryptionKey string // AES-256-GCM 密钥，须 32 字节，用于三要素与库内字段
	JWTSecret         string // 管理后台 JWT 签名密钥
}

// LogConfig 应用日志与接口审计日志保留策略。
type LogConfig struct {
	Level          string // zap 级别：debug/info/warn/error
	FilePath       string // 文件日志路径，空则仅 stdout
	RetentionDays  int    // api_request_logs 保留天数
	ArchiveEnabled bool   // 是否归档过期的文件日志
}

// AdminConfig 首次启动无管理员时创建的默认账号（生产环境务必修改）。
type AdminConfig struct {
	DefaultUsername string
	DefaultPassword string
}

// Load 从指定路径加载配置；文件不存在时使用默认值，环境变量 BRIDGE_* 可覆盖任意项。
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("BRIDGE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && !strings.Contains(err.Error(), "no such file") {
			return nil, err
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            v.GetInt("server.port"),
			Mode:            v.GetString("server.mode"),
			ReadTimeout:     v.GetDuration("server.read_timeout"),
			WriteTimeout:    v.GetDuration("server.write_timeout"),
			UpstreamTimeout: v.GetDuration("server.upstream_timeout"),
		},
		HuaAn: HuaAnConfig{
			BaseURL: v.GetString("huaan.base_url"),
			APIPath: v.GetString("huaan.api_path"),
		},
		Database: DatabaseConfig{
			DSN: v.GetString("database.dsn"),
		},
		Redis: RedisConfig{
			Addr:        v.GetString("redis.addr"),
			Password:    v.GetString("redis.password"),
			DB:          v.GetInt("redis.db"),
			BankListTTL: v.GetDuration("redis.bank_list_ttl"),
		},
		Security: SecurityConfig{
			DataEncryptionKey: v.GetString("security.data_encryption_key"),
			JWTSecret:         v.GetString("security.jwt_secret"),
		},
		Log: LogConfig{
			Level:          v.GetString("log.level"),
			FilePath:       v.GetString("log.file_path"),
			RetentionDays:  v.GetInt("log.retention_days"),
			ArchiveEnabled: v.GetBool("log.archive_enabled"),
		},
		Admin: AdminConfig{
			DefaultUsername: v.GetString("admin.default_username"),
			DefaultPassword: v.GetString("admin.default_password"),
		},
	}
	return cfg, nil
}

// setDefaults 为 viper 设置各配置项的默认值（配置文件缺失时生效）。
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.upstream_timeout", "25s")
	v.SetDefault("huaan.base_url", "https://xx.xxx.api.york.xin")
	v.SetDefault("huaan.api_path", "/upChannelApi")
	v.SetDefault("database.dsn", "bridge:bridge123@tcp(mysql:3306)/insurance_bridge?charset=utf8mb4&parseTime=True&loc=Local")
	v.SetDefault("redis.addr", "redis:6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.bank_list_ttl", "0")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.file_path", "./logs/app.log")
	v.SetDefault("log.retention_days", 90)
	v.SetDefault("log.archive_enabled", true)
	v.SetDefault("admin.default_username", "admin")
	v.SetDefault("admin.default_password", "admin123")
}
