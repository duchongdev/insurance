// Package job 后台定时任务：接口审计日志保留期清理、应用日志文件归档。
package job

import (
	"os"
	"path/filepath"
	"time"

	"github.com/huaan/insurance-bridge/internal/repository"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// CleanupJob 按 retentionDays 删除过期 api_request_logs，并可归档过期的 zap 文件日志。
type CleanupJob struct {
	logs           *repository.LogRepo
	log            *zap.Logger
	retentionDays  int    // 保留天数，如 90
	logFilePath    string // 应用日志文件路径
	archiveEnabled bool   // 是否将过期日志文件移动到 archive 子目录
}

// NewCleanupJob 构造清理任务，由 main 在启动后调用 Start。
func NewCleanupJob(logs *repository.LogRepo, log *zap.Logger, retentionDays int, logFilePath string, archiveEnabled bool) *CleanupJob {
	return &CleanupJob{
		logs:           logs,
		log:            log,
		retentionDays:  retentionDays,
		logFilePath:    logFilePath,
		archiveEnabled: archiveEnabled,
	}
}

// Start 注册 cron 表达式 "0 3 * * *"（每天 03:00）并启动调度器。
func (j *CleanupJob) Start() {
	c := cron.New()
	_, _ = c.AddFunc("0 3 * * *", j.run)
	c.Start()
}

// run 执行一次清理：删除早于 retentionDays 的 DB 日志，可选归档文件日志。
func (j *CleanupJob) run() {
	before := time.Now().AddDate(0, 0, -j.retentionDays)
	n, err := j.logs.DeleteBefore(before)
	if err != nil {
		j.log.Error("cleanup api logs failed", zap.Error(err))
	} else {
		j.log.Info("cleanup api logs done", zap.Int64("deleted", n))
	}
	if j.archiveEnabled && j.logFilePath != "" {
		j.archiveLogFile(before)
	}
}

// archiveLogFile 若当前日志文件修改时间早于 before，则重命名到 logs/archive/ 并创建新的空日志文件。
func (j *CleanupJob) archiveLogFile(before time.Time) {
	info, err := os.Stat(j.logFilePath)
	if err != nil || info.ModTime().After(before) {
		return
	}
	archiveDir := filepath.Join(filepath.Dir(j.logFilePath), "archive")
	_ = os.MkdirAll(archiveDir, 0o755)
	archiveName := filepath.Join(archiveDir, filepath.Base(j.logFilePath)+"."+before.Format("20060102"))
	if err := os.Rename(j.logFilePath, archiveName); err != nil {
		j.log.Warn("archive log file failed", zap.Error(err))
		return
	}
	f, err := os.OpenFile(j.logFilePath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err == nil {
		_ = f.Close()
	}
}
