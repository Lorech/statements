package commands

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statements",
		Short: "Bank statement parser",
		Long:  "Utility tool for automatically parsing and analyzing bank statements",
	}

	cmd.AddCommand(NewConfigCommand())
	cmd.AddCommand(NewProcessCommand())
	cmd.AddCommand(NewVersionCommand())

	return cmd
}
