package apisvr

import (
	"github.com/ahang7/go-IAM/internal/pkg/code"
	"github.com/ahang7/go-IAM/internal/pkg/middleware/auth"
	httpcore "github.com/ahang7/go-IAM/pkg/core/http"
	"github.com/ahang7/go-sdk/errors"
	"github.com/gin-gonic/gin"
)

func router(g *gin.Engine) {
	installMiddleware(g)
	restController(g)
}

func installMiddleware(g *gin.Engine) {}

func restController(g *gin.Engine) {
	// Middlewares
	strategy := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", strategy.LoginHandler)
	g.POST("/logout", strategy.LogoutHandler)
	g.POST("/refresh", strategy.RefreshHandler)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthExecute(), func(ctx *gin.Context) {
		httpcore.WriteResponse(ctx,
			errors.WithCode(code.ErrPageNotFound, "page not found"),
			nil,
		)
	})
}
