package auth

import (
	"strings"

	"github.com/ahang7/go-IAM/internal/pkg/code"
	"github.com/ahang7/go-IAM/internal/pkg/middleware"
	httpcore "github.com/ahang7/go-IAM/pkg/core/http"
	"github.com/ahang7/go-IAM/pkg/errors"
	"github.com/gin-gonic/gin"
)

type AutoStrategy struct {
	basic middleware.AuthStrategy
	jwt   middleware.AuthStrategy
}

func NewAutoStrategy(basic, jwt middleware.AuthStrategy) AutoStrategy {
	return AutoStrategy{
		basic: basic,
		jwt:   jwt,
	}
}

func (a AutoStrategy) SetBasicStrategy(basic middleware.AuthStrategy) {
	a.basic = basic
}

func (a AutoStrategy) SetJwtStrategy(jwt middleware.AuthStrategy) {
	a.jwt = jwt
}

func (a AutoStrategy) AuthExecute() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		operator := middleware.AuthOperator{}

		if len(authHeader) != 2 {
			httpcore.WriteResponse(c,
				errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()
			return
		}

		switch authHeader[0] {
		case "Basic":
			operator.SetStrategy(a.basic)
		case "Bearer":
			operator.SetStrategy(a.jwt)
		default:
			httpcore.WriteResponse(c,
				errors.WithCode(code.ErrSignatureInvalid, "unrecognized auth type"),
				nil,
			)
			c.Abort()
			return
		}

		operator.AuthExecute()(c)

		c.Next()
	}
}
