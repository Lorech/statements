package adapters

import (
	"statements/pkg/config"
	"statements/pkg/transactions"
)

// A shared interface for bank transaction structs.
type TransactionAdapter interface {
	FieldValue(field string) any
	Normalize() transactions.Transaction
}

// Converts a slice of adapter-conforming items into a slice of adapter interface items.
func AdaptTransactions[T TransactionAdapter](ts []T) []TransactionAdapter {
	adapters := make([]TransactionAdapter, len(ts))
	for i, v := range ts {
		adapters[i] = v
	}
	return adapters
}

// Filters a slice of transactions based on the configured filters.
func FilterTransactions[T TransactionAdapter](ts []T, fs []config.Filter) []TransactionAdapter {
	var res []TransactionAdapter
	for _, t := range ts {
		valid := true
		for _, f := range fs {
			if !f.Match(t.FieldValue(f.FieldName())) {
				valid = false
				break
			}
		}
		if valid {
			res = append(res, t)
		}
	}
	return res
}
