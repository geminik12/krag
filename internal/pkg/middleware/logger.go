package middleware

import (
	"log/slog"
	"time"

	"github.com/geminik12/krag/internal/pkg/contextx"
	"github.com/geminik12/krag/internal/pkg/known"
	"github.com/gin-gonic/gin"
)

// Logger 记录 HTTP 请求的详细信息.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		cost := time.Since(start)

		// 获取请求相关信息
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		status := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		// 从上下文中获取 RequestID，确保整个链路的日志都能关联起来
		requestID := c.Request.Header.Get(known.XRequestID)
		if requestID == "" {
			requestID = contextx.RequestID(c.Request.Context())
		}

		// 构造日志字段
		fields := []any{
			slog.String("request_id", requestID),
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", ip),
			slog.String("user_agent", userAgent),
			// 将耗时单位统一为毫秒，方便查看
			slog.Float64("cost_ms", float64(cost.Microseconds())/1000.0),
		}

		if errorMessage != "" {
			fields = append(fields, slog.String("error", errorMessage))
		}

		// 根据状态码记录不同级别的日志
		if status >= 500 {
			slog.Error("HTTP Request", fields...)
		} else if status >= 400 {
			slog.Warn("HTTP Request", fields...)
		} else {
			slog.Info("HTTP Request", fields...)
		}
	}
}
