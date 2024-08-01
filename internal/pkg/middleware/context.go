package middleware

import "github.com/gin-gonic/gin"

const (
	UserNameKey = "username"
)

func Context() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("requestID", ctx.GetString(XRequestIDKey))
		ctx.Set("username", ctx.GetString(UserNameKey))
		ctx.Next()
	}
}
