package commands

import (
	"encoding/csv"
	"fmt"
	"os"

	"statements/pkg/adapters"
	"statements/pkg/config"
	"statements/pkg/transactions"

	"github.com/spf13/cobra"
)

func NewProcessCommand() *cobra.Command {
	var infile, outfile, confile *string

	cmd := &cobra.Command{
		Use:   "process",
		Short: "Process a bank statement",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Parse(*confile)
			if err != nil {
				return fmt.Errorf("could not parse config file: %v", err)
			}

			var bank transactions.Bank
			err = bank.Set(c.Flags.Bank)
			if err != nil {
				return err
			}

			if *infile == "" {
				cInput := c.Flags.Input
				if cInput != "" {
					*infile = cInput
				} else {
					bInput, err := bank.Input()
					if err == nil {
						*infile = bInput
					} else {
						return fmt.Errorf("no input file provided")
					}
				}
			}

			records, err := readInput(*infile)
			if err != nil {
				return err
			}

			var bts []adapters.TransactionAdapter
			var fs []config.Filter

			switch bank {
			case "swedbank":
				sts, err := adapters.NewSwedbankTransactions(records)
				if err != nil {
					return err
				}
				bts = adapters.AdaptTransactions(sts)

				for _, rf := range c.Filters {
					f, err := rf.DecodeWithFieldMap(adapters.SwedbankFieldMap)
					if err != nil {
						return err
					}
					fs = append(fs, f)
				}
			}

			bts = adapters.FilterTransactions(bts, fs)

			var ts []transactions.Transaction
			for _, bt := range bts {
				ts = append(ts, bt.Normalize())
			}

			// TODO: Add classification here

			if *outfile == "" {
				cOutput := c.Flags.Output
				if cOutput != "" {
					*outfile = cOutput
				} else {
					*outfile = "output.csv"
				}
			}

			err = writeOutput(*outfile, ts)
			if err != nil {
				return err
			}

			return nil
		},
	}

	infile = cmd.Flags().StringP("input", "i", "", "input file to process")
	cmd.MarkFlagFilename("input")

	outfile = cmd.Flags().StringP("output", "o", "", "output file to write to")

	confile = cmd.Flags().String("config", config.DefaultConfig, "configuration file to use")

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
