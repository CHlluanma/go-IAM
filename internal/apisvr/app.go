package apisvr

import (
	"github.com/ahang7/go-IAM/internal/apisvr/options"
	"github.com/ahang7/go-IAM/pkg/app"
)

const commandDesc = `The IAM API server validates and configures data ...`

func NewApp() *app.App {
	// TODO: 配置appoptions
	opts := options.NewOptions()
	a := app.NewApp("IAM",
		"iamsvr",
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)
	return a
}

func run(opts *options.Options) app.RunFunc {
	return func(app string) error {
		// TODO: 初始化日志

		//TODO: 读取配置

		return nil
	}
}
