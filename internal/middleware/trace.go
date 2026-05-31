// Package middleware 提供 Gin 全局中间件。
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceHeader 请求/响应头中的链路 ID 字段名，可与审计日志 trace_id 关联（代理层另生成 UUID 写库）。
const TraceHeader = "X-Trace-Id"

// Trace 中间件：从请求头读取或生成 traceId，写入 gin.Context 并回写响应头。
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(TraceHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("traceId", traceID)
		c.Writer.Header().Set(TraceHeader, traceID)
		c.Next()
	}
}
