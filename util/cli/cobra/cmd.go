package cobra

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CmdOption func(cmd *cobra.Command)

func NewCmd(name string, shortDesc string, example string, options ...CmdOption) *cobra.Command {
	cmd := cobra.Command{
		Use:       name,
		Short:     shortDesc,
		Example:   "  " + name + " " + example,
		ValidArgs: []string{},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.IsAvailableCommand() {
				// don't try to load the config if an unavailable command is being run
				//   (this allows the "help" command to execute without error)
				return nil
			}

			// don't show help if the command produces an error
			cmd.SilenceUsage = true

			return nil
		},
	}

	for _, opt := range options {
		opt(&cmd)
	}

	return &cmd
}

func WithFlags(flags *pflag.FlagSet) CmdOption {
	return func(cmd *cobra.Command) {
		if err := cmd.ParseFlags(flags.Args()); err != nil {
			panic(err.Error())
		}
	}
}

func WithArgs(args cobra.PositionalArgs) CmdOption {
	return func(cmd *cobra.Command) {
		cmd.Args = args
	}
}

func WithRunE(fn func(cmd *cobra.Command, args []string) error) CmdOption {
	return func(cmd *cobra.Command) {
		cmd.RunE = fn
	}
}
