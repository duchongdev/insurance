// Package logger 初始化 zap 日志：JSON 格式输出到 stdout 与可选轮转文件（lumberjack）。
package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Options 日志初始化参数。
type Options struct {
	Level          string // debug / info / warn / error
	FilePath       string // 空则仅 stdout
	RetentionDays  int    // 轮转文件保留天数，超期自动删除
	ArchiveEnabled bool   // 轮转后 gzip 压缩
	MaxSizeMB      int    // 单文件上限（MB），达到后轮转
}

// New 创建生产级 JSON 日志器；filePath 非空时使用 lumberjack 轮转、压缩与按天龄清理。
func New(opts Options) (*zap.Logger, error) {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(opts.Level)); err != nil {
		lvl.SetLevel(zap.InfoLevel)
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "timestamp"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.MessageKey = "message"
	encCfg.LevelKey = "level"
	encCfg.CallerKey = "caller"

	encoder := zapcore.NewJSONEncoder(encCfg)
	cores := []zapcore.Core{zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl)}

	if opts.FilePath != "" {
		if err := os.MkdirAll(filepath.Dir(opts.FilePath), 0o755); err != nil {
			return nil, err
		}
		maxSize := opts.MaxSizeMB
		if maxSize <= 0 {
			maxSize = 100
		}
		retention := opts.RetentionDays
		if retention <= 0 {
			retention = 90
		}
		lw := &lumberjack.Logger{
			Filename:   opts.FilePath,
			MaxSize:    maxSize,
			MaxAge:     retention,
			MaxBackups: 0,
			Compress:   opts.ArchiveEnabled,
			LocalTime:  true,
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(lw), lvl))
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)), nil
}
