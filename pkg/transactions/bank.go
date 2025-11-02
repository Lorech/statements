package transactions

import (
	"fmt"
	"strings"
)

type Bank string

const (
	BankSwedbank Bank = "swedbank"
)

// The stringified representation of the bank enumeration.
//
// Used as help text in Cobra.
func (b *Bank) Type() string {
	return "Bank"
}

// The stringified representation of a bank.
//
// Used for printing and as help text in Cobra.
func (b *Bank) String() string {
	return string(*b)
}

// Set the value of a bank.
//
// Also used as a way to convert `string` to `Bank`
func (b *Bank) Set(v string) error {
	vp := strings.ToLower(v)
	switch vp {
	case "swedbank":
		*b = Bank(vp)
		return nil
	default:
		return fmt.Errorf(`must be one of "swedbank"`)
	}
}

// The default input file for a bank.
func (b *Bank) Input() (string, error) {
	switch *b {
	case BankSwedbank:
		return "statement.csv", nil
	default:
		return "", fmt.Errorf(`no default input file specified`)
	}
}
