package commands

import (
	"fmt"
	"statements/pkg/config"

	"github.com/spf13/cobra"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure the CLI",
	}

	validateCmd := &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate a configuration file, defaulting to config.json",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := config.DefaultConfig
			if len(args) == 1 {
				configFile = args[0]
			}

			err := config.Validate(configFile)
			if err != nil {
				return err
			}

			fmt.Printf("Configuration valid!\n")
			return nil
		},
	}

	validateCmd.DisableFlagsInUseLine = true

	cmd.AddCommand(validateCmd)

	return cmd
}
