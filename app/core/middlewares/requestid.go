package middlewares

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

const (
	XRequestIDKey = "X-Request-ID"
)

// generator a function type that returns string.
type generator func() string

var (
	random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
)

func uuid(len int) string {
	bytes := make([]byte, len)
	random.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)[:len]
}

// RequestID is a middleware that injects a 'RequestID' into the context and header of each request.
func RequestID(gen generator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var xRequestID string
		if gen != nil {
			xRequestID = gen()
		} else {
			xRequestID = uuid(16)
		}
		c.Request.Header.Set(XRequestIDKey, xRequestID)
		c.Set(XRequestIDKey, xRequestID)
		c.Next()
	}
}
