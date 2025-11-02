package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the tool's version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("statements v0.1.0")
		},
	}
}
