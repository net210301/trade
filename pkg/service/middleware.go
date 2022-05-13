package service

import (
	"github.com/gin-gonic/gin"
	"github.com/net210301/trade/pkg/tlog"
	"go.uber.org/zap"
	"net/http"
)

func recoverMiddleware(logger *tlog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var fields = make([]tlog.Field, 0, 8)
		defer func() {
			if rec := recover(); rec != nil {
				var err = rec.(error)
				fields = append(fields, zap.String("err", err.Error()))
				logger.Error("access", fields...)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			fields = append(fields,
				zap.String("method", c.Request.Method),
				zap.Int("code", c.Writer.Status()),
				zap.Int("size", c.Writer.Size()),
				zap.String("host", c.Request.Host),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.String("agent", c.Request.UserAgent()),
				zap.String("proto", c.Request.Proto),
				zap.String("err", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)
			logger.Info("access", fields...)
		}()
		c.Next()
	}
}
