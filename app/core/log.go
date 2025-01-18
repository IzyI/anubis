package core

import (
	"anubis/app/core/middlewares"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// GetLoggerConfig return gin.LoggerConfig which will write the logs to specified io.Writer with given gin.LogFormatter.
// By default gin.DefaultWriter = os.Stdout
// reference: https://github.com/gin-gonic/gin#custom-log-format
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

// GetDefaultLogFormatterWithRequestID returns gin.LogFormatter with 'RequestID'
func GetDefaultLogFormatterWithRequestID() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		return fmt.Sprintf(
			"[GIN-debug] %s | %s | %s | %s | %s | %3d | %s | %s | %s\n",
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

func LogDebag(ctx *gin.Context, s string) {
	fmt.Printf("[GIN-debug] %s | %s | %s | %s | %s | %s\n",
		ctx.Request.Method,
		time.Now().UTC(),
		ctx.Request.Header.Get(middlewares.XRequestIDKey),
		ctx.Request.URL,
		ctx.ClientIP,
		s)
}
