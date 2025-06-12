package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

type RequestLogger struct {
	log *zerolog.Logger
}

func NewRequestLogger(log *zerolog.Logger) *RequestLogger {
	return &RequestLogger{log}
}

func (l *RequestLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		logEvent := l.log.Info()

		if status >= 400 {
			logEvent = l.log.Warn()
		}
		if status >= 500 {
			logEvent = l.log.Error()
		}

		logEvent.
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("duration", duration).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg(fmt.Sprintf("%s %s %d %s", method, path, status, duration))
	}
}
