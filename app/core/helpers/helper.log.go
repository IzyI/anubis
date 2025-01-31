package helpers

import (
	"anubis/app/core/middlewares"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

func GetLoggerConfig(formatter gin.LogFormatter, output io.Writer, skipPaths []string) gin.LoggerConfig {
	if formatter == nil {
		formatter = GetDefaultLogFormatterWithRequestID()
	}
	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skipPaths,
	}
}

func GetDefaultLogFormatterWithRequestID() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		return fmt.Sprintf(
			"[GIN] %s | %s | %s | %s | %s | %3d | %s | %s | %s",
			param.Method,
			param.TimeStamp.Format(time.RFC3339),
			param.Request.Header.Get(middlewares.XRequestIDKey),
			param.Path,
			param.ClientIP,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}
}

func LogDebug(ctx *gin.Context, s string) {
	fmt.Printf("[GIN-error] %s | %s | %s | %s | %s | %s",
		ctx.Request.Method,
		time.Now().UTC(),
		ctx.Request.Header.Get(middlewares.XRequestIDKey),
		ctx.Request.URL,
		ctx.ClientIP,
		s)
}
