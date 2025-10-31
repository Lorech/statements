package commands

import (
	"errors"
	"fmt"

	"statements/pkg/transactions"

	"github.com/spf13/cobra"
)

func NewProcessCommand() *cobra.Command {
	var bank transactions.Bank
	var input, output *string

	cmd := &cobra.Command{
		Use:   "process",
		Short: "Process a bank statement",
		RunE: func(cmd *cobra.Command, args []string) error {
			if *input == "" {
				bInput, err := bank.Input()
				if err == nil {
					*input = bInput
				} else {
					return errors.New("bank has no default input, provide one manually")
				}
			}

			fmt.Printf("Bank: %s; Input: %s; Output: %s", bank, *input, *output)

			return nil
		},
	}

	cmd.Flags().VarP(&bank, "bank", "b", `bank to parse input as, options - "swedbank"`)
	cmd.MarkFlagRequired("bank")

	input = cmd.Flags().StringP("input", "i", "", "input file to process")
	cmd.MarkFlagFilename("input")

	output = cmd.Flags().StringP("output", "o", "output.csv", "output file to write to")

	return cmd
}
