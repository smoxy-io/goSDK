package cli

import "github.com/spf13/pflag"

func ParseArgs(args []string, opts ...FlagOption) (*pflag.FlagSet, error) {
	flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)

	if len(args) == 0 {
		return flagSet, nil
	}

	if len(opts) != 0 {
		for _, opt := range opts {
			opt(flagSet)
		}
	}

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	return flagSet, nil
}
