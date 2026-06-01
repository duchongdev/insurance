// Package job 后台定时任务：清理超过保留期的历史日志文件。
package job

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// CleanupJob 扫描日志目录，删除修改时间早于 retentionDays 的日志及压缩归档文件。
type CleanupJob struct {
	log           *zap.Logger
	retentionDays int
	logFilePath   string
}

// NewCleanupJob 构造清理任务，由 main 在启动后调用 Start。
func NewCleanupJob(log *zap.Logger, retentionDays int, logFilePath string) *CleanupJob {
	return &CleanupJob{
		log:           log,
		retentionDays: retentionDays,
		logFilePath:   logFilePath,
	}
}

// Start 注册 cron 表达式 "0 3 * * *"（每天 03:00）并启动调度器。
func (j *CleanupJob) Start() {
	if j.logFilePath == "" {
		return
	}
	c := cron.New()
	_, _ = c.AddFunc("0 3 * * *", j.run)
	c.Start()
}

// run 删除 logs 目录及 archive 子目录中超过保留期的 .log / .gz 文件。
func (j *CleanupJob) run() {
	retention := j.retentionDays
	if retention <= 0 {
		retention = 90
	}
	before := time.Now().AddDate(0, 0, -retention)
	logDir := filepath.Dir(j.logFilePath)

	var removed int
	_ = filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		name := info.Name()
		if !strings.HasSuffix(name, ".log") && !strings.HasSuffix(name, ".gz") {
			return nil
		}
		// 当前正在写入的主日志由 lumberjack 管理，跳过
		if path == j.logFilePath {
			return nil
		}
		if info.ModTime().Before(before) {
			if err := os.Remove(path); err == nil {
				removed++
			}
		}
		return nil
	})

	j.log.Info("cleanup log files done",
		zap.String("logDir", logDir),
		zap.Int("retentionDays", retention),
		zap.Int("removed", removed),
	)
}
