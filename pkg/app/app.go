package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/ahang7/go-IAM/pkg/log"
)

type App struct {
	prefix      string
	appname     string
	description string

	runFunc RunFunc
	flags   FlagsOptions

	noConfig  bool
	noVersion bool

	commands []*Command
	use      string
	short    string
	long     string
	args     cobra.PositionalArgs
	cmd      *cobra.Command
}

type RunFunc func(app string) error

type Option func(*App)

func WithFlags(fi FlagsOptions) Option {
	return func(app *App) {
		app.flags = fi
	}
}

func WithDescription(desc string) Option {
	return func(app *App) {
		app.description = desc
	}
}

func WithNoConfig() Option {
	return func(app *App) {
		app.noConfig = true
	}
}

func WithVersion(noVersion bool) Option {
	return func(app *App) {
		app.noVersion = noVersion
	}
}

func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runFunc = run
	}
}

// WithCommand 设置命令行参数
func WithCommand(use, short, long string) Option {
	return func(app *App) {
		app.use = use
		app.short = short
		app.long = long
	}
}

// WithDefaultValidArgs set default valid args to valid non-flag arguments
func WithDefaultValidArgs() Option {
	return func(app *App) {
		app.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("unknown command %q", arg)
				}
			}
			return nil
		}
	}
}

func NewApp(prefix, appname string, opts ...Option) *App {
	a := &App{
		prefix:  prefix,
		appname: appname,
	}

	for _, o := range opts {
		o(a)
	}
	a.validateCommand()
	a.buildCommand()

	return a
}

func (a *App) validateCommand() {
	if a.use == "" {
		a.use = a.appname
	}
	if a.long == "" {
		a.long = a.description
	}
}

func (a *App) buildCommand() {
	cmd := &cobra.Command{
		Use:   a.use,
		Short: a.short,
		Long:  a.long,
		Args:  a.args,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SetNormalizeFunc(normalizeFunc)

	if len(a.commands) > 0 {
		for _, c := range a.commands {
			cmd.AddCommand(c.buildCommand())
		}
		cmd.SetHelpCommand(helpCommand(a.use))
	}
	if a.runFunc != nil {
		cmd.RunE = a.run
	}

	// 构建命令行参数解析
	var appFlags FlagSet
	if a.flags != nil {
		appFlags = a.flags.Flags()
		fs := cmd.Flags()
		for _, f := range appFlags.flags {
			fs.AddFlagSet(f)
		}
	}
	// config 命令行
	if !a.noConfig {
		addConfigFile(a.prefix, a.appname, appFlags.Flags("config"))
	}
	// version 命令行
	if !a.noVersion {
		addConfigFile(a.prefix, a.appname, appFlags.Flags("version"))
	}
	// 全局 命令行
	appFlags.Flags("global").BoolP("help", "h", false, fmt.Sprintf("help for %s", color.GreenString(a.appname)))
	cmd.Flags().AddFlagSet(appFlags.Flags("global"))

	a.cmd = cmd
}

func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		log.Fatal("app run failed")
	}
}

func (a *App) run(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		log.Debugf("flag %s: %v", flag.Name, flag.Value)
	})
	// todo: 输出 --version

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viper.Unmarshal(a.flags); err != nil {
			return err
		}
	}

	if a.runFunc != nil {
		return a.runFunc(a.appname)
	}
	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v working dir: %s", color.GreenString("===>"), wd)
}
