// Package logger 初始化 zap 日志：JSON 格式同时输出到 stdout 与可选文件。
package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New 创建生产级 JSON 日志器；level 解析失败时默认为 info；filePath 非空时追加文件输出。
func New(level, filePath string) (*zap.Logger, error) {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
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

	if filePath != "" {
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return nil, err
		}
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(f), lvl))
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)), nil
}
