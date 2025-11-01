package commands

import (
	"encoding/csv"
	"fmt"
	"os"

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
					return fmt.Errorf("bank has no default input, provide one manually")
				}
			}

			records, err := readInput(*input)
			if err != nil {
				return err
			}

			err = writeOutput(*output, records)
			if err != nil {
				return err
			}

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

// Reads the provided input file, returning a parsed CSV if successful.
func readInput(input string) ([][]string, error) {
	records := [][]string{}

	in, err := os.Open(input)
	if err != nil {
		return records, fmt.Errorf("input file could not be opened: %v", err)
	}
	defer in.Close()

	r := csv.NewReader(in)
	r.Comma = ';'
	records, err = r.ReadAll()
	if err != nil {
		return records, fmt.Errorf("csv file could not be read: %v", err)
	}

	return records, nil
}

// Writes the provided CSV data to the provided output file.
func writeOutput(output string, data [][]string) error {
	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("output file could not be opened: %v", err)
	}
	defer out.Close()

	w := csv.NewWriter(out)
	w.Comma = ';'
	err = w.WriteAll(data)
	if err != nil {
		return fmt.Errorf("csv file could not be written: %v", err)
	}

	return nil
}
