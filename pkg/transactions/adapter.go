package transactions

// A shared interface for bank transaction structs.
type TransactionAdapter interface {
	Normalize() Transaction
}

// Converts a slice of TransactionAdapter-conforming items into a slice of adapter values.
func AdaptTransactions[T TransactionAdapter](ts []T) []TransactionAdapter {
	adapters := make([]TransactionAdapter, len(ts))
	for i, v := range ts {
		adapters[i] = v
	}
	return adapters
}
