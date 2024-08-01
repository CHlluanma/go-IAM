package httpcore

import (
	"github.com/ahang7/go-IAM/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Response 通用返回包
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ErrResponse 错误返回包
type ErrResponse struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Reference string `json:"reference,omitempty"`
}

// WriteResponse 返回http response
func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), &ErrResponse{
			Code:      coder.HTTPStatus(),
			Msg:       coder.String(),
			Reference: coder.Reference(),
		})
	}
	c.JSON(200, data)
}
