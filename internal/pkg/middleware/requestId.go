package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	XRequestIDKey = "X-Request-ID"
)

// RequestID 返回一个gin中间件函数，用于生成和注入请求ID。
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rid := ctx.GetHeader(XRequestIDKey)
		if rid == "" {
			rid = uuid.Must(uuid.NewV7()).String()
			ctx.Request.Header.Set(XRequestIDKey, rid)
			ctx.Set(XRequestIDKey, rid)
		}
		ctx.Writer.Header().Set(XRequestIDKey, rid)
		ctx.Next()
	}
}
