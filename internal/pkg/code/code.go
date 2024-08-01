package code

import (
	"github.com/ahang7/go-IAM/pkg/errors"
	"github.com/novalagung/gubrak"
)

type ErrCode struct {
	C    int
	HTTP int
	Ext  string
	Ref  string
}

// Code implements errors.Coder.
func (c *ErrCode) Code() int {
	return c.HTTP
}

// HTTPStatus implements errors.Coder.
func (c *ErrCode) HTTPStatus() int {
	return c.HTTP

}

// Reference implements errors.Coder.
func (c *ErrCode) Reference() string {
	return c.Ref

}

// String implements errors.Coder.
func (c *ErrCode) String() string {
	return c.Ext

}

var _ errors.Coder = (*ErrCode)(nil)

func register(code int, httpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes([]int{200, 400, 401, 404, 500}, httpStatus)
	if !found {
		panic("http status code must be 200, 400, 401, 404, 500")
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	errors.MustRegister(&ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	})
}
