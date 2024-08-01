package auth

import (
	"encoding/base64"
	"strings"

	"github.com/ahang7/go-IAM/internal/pkg/code"
	"github.com/ahang7/go-IAM/internal/pkg/middleware"
	httpcore "github.com/ahang7/go-IAM/pkg/core/http"
	"github.com/ahang7/go-IAM/pkg/errors"
	"github.com/gin-gonic/gin"
)

// BasicStrategy basic auth strategy
type BasicStrategy struct {
	compare basicCompare
}

// NewBasicStrategy create basic strategy
func NewBasicStrategy(compare basicCompare) BasicStrategy {
	return BasicStrategy{compare: compare}
}

type basicCompare func(username, password string) bool

func (b BasicStrategy) AuthExecute() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			httpcore.WriteResponse(c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil)
			c.Abort()
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || !b.compare(pair[0], pair[1]) {
			httpcore.WriteResponse(c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil)
			c.Abort()
			return
		}

		c.Set(middleware.UserNameKey, pair[0])
		c.Next()
	}
}
