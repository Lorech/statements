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

			var bts []transactions.TransactionAdapter
			switch bank {
			case "swedbank":
				sts, err := transactions.NewSwedbankTransactions(records)
				if err != nil {
					return err
				}
				bts = transactions.AdaptTransactions(sts)
			}

			// TODO: Add filtering here

			var ts []transactions.Transaction
			for _, bt := range bts {
				ts = append(ts, bt.Normalize())
			}

			// TODO: Add classification here

			err = writeOutput(*output, ts)
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

// Writes transactions to the provided output file as CSV.
func writeOutput(output string, transactions []transactions.Transaction) error {
	var rows [][]string
	for _, t := range transactions {
		rows = append(rows, t.Csv())
	}

	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("output file could not be opened: %v", err)
	}
	defer out.Close()

	w := csv.NewWriter(out)
	w.Comma = ';'
	err = w.WriteAll(rows)
	if err != nil {
		return fmt.Errorf("csv file could not be written: %v", err)
	}

	return nil
}
