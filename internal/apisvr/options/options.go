package options

import (
	"encoding/json"

	pkgoptions "github.com/ahang7/go-IAM/internal/pkg/options"
	"github.com/ahang7/go-IAM/pkg/app"
)

type Options struct {
	GenericServerRunOptions *pkgoptions.ServerRunOptions `json:"server" mapstructure:"server"`
	MySQLOpts               *pkgoptions.MySQLOptions     `json:"mysql" mapstructure:"mysql"`
}

func (o *Options) Complete() error {

	return nil
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func (o *Options) ApplyFlags() []error {
	return nil
}

func (o *Options) Flags() (fs app.FlagSet) {
	o.MySQLOpts.AddFlags(fs.Flags("mysql"))

	return
}

var _ app.OptionsIntf = (*Options)(nil)

func NewOptions() *Options {
	o := &Options{
		MySQLOpts: pkgoptions.NewMySQLOptionsNil(),
	}
	return o
}
