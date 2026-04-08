package middleware

import (
	"time"

	"github.com/Knowpals/Knowpals-be-go/pkg/otelx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	logger *zap.Logger
}

func NewLoggerMiddleware(l *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger: l}
}

func (l *LoggerMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		traceID := otelx.TraceIDFromContext(ctx.Request.Context())
		spanID := otelx.SpanIDFromContext(ctx.Request.Context())
		//处理请求
		ctx.Next()

		cost := time.Since(start)
		fields := []zap.Field{
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.Int("status", ctx.Writer.Status()),
			zap.String("client_ip", ctx.ClientIP()),
			zap.Duration("latency", cost),
		}
		if traceID != "" {
			fields = append(fields, zap.String("trace_id", traceID))
		}
		if spanID != "" {
			fields = append(fields, zap.String("span_id", spanID))
		}
		if len(ctx.Errors) > 0 {
			fields = append(fields,
				zap.String("errors", ctx.Errors.String()),
			)
			l.logger.Error("HTTP request error", fields...)
		} else {
			l.logger.Info("HTTP request success", fields...)
		}
	}
}
