package auth

import (
	"github.com/ahang7/go-IAM/internal/pkg/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// AuthzAudience defines the audience of the token
const AuthzAudience = "iam.authz.ch.com"

type JWTStrategy struct {
	jwt.GinJWTMiddleware
}

var _ middleware.AuthStrategy = &JWTStrategy{}

// NewJWTStrategy creates a new JWT strategy
func NewJWTStrategy(gjwt jwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{gjwt}
}

func (j JWTStrategy) AuthExecute() gin.HandlerFunc {
	return j.MiddlewareFunc()
}
