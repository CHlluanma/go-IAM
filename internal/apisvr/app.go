package apisvr

import (
	"github.com/ahang7/go-IAM/internal/apisvr/options"
	"github.com/ahang7/go-IAM/pkg/app"
	"github.com/ahang7/go-IAM/pkg/log"
)

const commandDesc = `The IAM API server validates and configures data ...`

func NewApp() *app.App {
	// 配置appoptions
	opts := options.NewOptions()
	a := app.NewApp(
		"IAM",
		"iamsvr",
		app.WithFlags(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)
	return a
}

func run(opts *options.Options) app.RunFunc {
	return func(app string) error {
		log.Infof("opts: %v", opts)
		return nil
	}
}
