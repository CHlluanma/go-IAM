package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

type Command struct {
	usage string
	short string
	long  string

	flags FlagIntf

	subCommand []*Command

	// args 用于验证参数合法性
	args    cobra.PositionalArgs
	runFunc RunCommandFunc
}

type RunCommandFunc func(args []string) error

func (c *Command) AddCommand(cmd *Command) {
	c.subCommand = append(c.subCommand, cmd)
}

func (c *Command) AddCommands(cmds ...*Command) {
	c.subCommand = append(c.subCommand, cmds...)
}

func (c *Command) buildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.short,
		Long:  c.long,
		Args:  c.args,
	}
	cmd.SetOut(os.Stdout)

	// 添加子命令
	if len(c.subCommand) > 0 {
		for _, command := range c.subCommand {
			cmd.AddCommand(command.buildCommand())
		}
	}
	if c.runFunc != nil {
		cmd.Run = c.run
	}

	// 添加命令行
	if c.flags != nil {
		for _, f := range c.flags.Flags().flags {
			cmd.Flags().AddFlagSet(f)
		}
	}

	// 添加Help命令
	addHelpCommandFlag(c.usage, cmd.Flags())
	// 添加Version命令

	return cmd
}

func (c *Command) run(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			log.Fatal("run command function failed")
		}
	}
}
