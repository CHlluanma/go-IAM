package app

import (
	"github.com/spf13/pflag"
)

type FlagSet struct {
	flags map[string]*pflag.FlagSet
}

func (fs *FlagSet) Flags(name string) *pflag.FlagSet {
	if fs.flags == nil {
		fs.flags = make(map[string]*pflag.FlagSet)
	}
	if _, ok := fs.flags[name]; !ok {
		fs.flags[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
	}
	return fs.flags[name]
}
