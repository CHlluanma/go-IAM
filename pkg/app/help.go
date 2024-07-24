package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

const helpTemplate = `
Usage:
	{{.UseLine}}
Description:
	{{.Short}}
Commands:
{{- range .Commands}}
	{{.Name}}: {{if .Short}}
{{- end}}
`

func helpCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "help about any command",
		Long: `Help provides help for any command in the application.
		Simply type ` + name + `help [path to command] for full details`,
		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknown help topic %#q\n", args[0])
			} else {
				cmd.InitDefaultHelpFlag()
				_ = cmd.Help()
			}
		},
	}
}

func addHelpCommandFlag(use string, fs *pflag.FlagSet) {
	fs.BoolP(
		"help",
		"H",
		false,
		fmt.Sprintf("help for the %s command.", color.GreenString(strings.Split(use, " ")[0])),
	)
}
