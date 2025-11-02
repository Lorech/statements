package transactions

import (
	"fmt"
	"strings"
)

type Bank string

const (
	BankSwedbank Bank = "swedbank"
)

// The stringified representation of a bank.
func (b *Bank) String() string {
	return string(*b)
}

// Parses a string into a bank.
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
