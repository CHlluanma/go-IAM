package options

import (
	pkgoptions "github.com/ahang7/go-IAM/internal/pkg/options"
	"github.com/ahang7/go-IAM/pkg/app"
)

type Options struct {
	MySQLOpts *pkgoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
}

func (o *Options) Complete() error {
	//TODO implement me
	panic("implement me")
}

func (o *Options) String() string {
	//TODO implement me
	panic("implement me")
}

func (o *Options) ApplyFlags() []error {
	//TODO implement me
	panic("implement me")
}

func (o *Options) Flags() (fs app.FlagSet) {
	//TODO implement me
	panic("implement me")
}

func NewOptions() *Options {
	o := &Options{
		MySQLOpts: pkgoptions.NewMySQLOptionsNil(),
	}
	return o
}
