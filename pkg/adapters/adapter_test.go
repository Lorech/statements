package adapters_test

import (
	"statements/pkg/adapters"
	"statements/pkg/config"
	"statements/pkg/transactions"
	"testing"
	"time"
)

// Mock transaction adapter for testing
type mockTransaction struct {
	date  time.Time
	value uint
	desc  string
}

func (m mockTransaction) FieldValue(field string) any {
	switch field {
	case "date":
		return m.date
	case "value":
		return float64(m.value)
	case "description":
		return m.desc
	}
	return nil
}

func (m mockTransaction) Normalize() transactions.Transaction {
	return transactions.Transaction{
		Date:          m.date,
		AccountHolder: "Test",
		Description:   m.desc,
		Value:         int(m.value),
		Currency:      "EUR",
	}
}

func TestAdaptTransactions(t *testing.T) {
	mockTxs := []mockTransaction{
		{
			date:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			value: 100,
			desc:  "Transaction 1",
		},
		{
			date:  time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			value: 200,
			desc:  "Transaction 2",
		},
	}

	adapted := adapters.AdaptTransactions(mockTxs)

	if len(adapted) != len(mockTxs) {
		t.Errorf("AdaptTransactions() returned %d items, want %d", len(adapted), len(mockTxs))
	}

	for i, tx := range adapted {
		if tx == nil {
			t.Errorf("AdaptTransactions() item %d is nil", i)
		}
	}
}

func TestFilterTransactions(t *testing.T) {
	mockTxs := []mockTransaction{
		{
			date:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			value: 100,
			desc:  "Test transaction",
		},
		{
			date:  time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			value: 200,
			desc:  "Another transaction",
		},
		{
			date:  time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			value: 300,
			desc:  "Third transaction",
		},
	}

	tests := []struct {
		name        string
		filters     []config.Filter
		wantCount   int
		description string
	}{
		{
			name:        "no filters - all pass",
			filters:     []config.Filter{},
			wantCount:   3,
			description: "All transactions should pass with no filters",
		},
		{
			name: "date filter - equal",
			filters: []config.Filter{
				config.DateFilter{
					Field:      "date",
					Condition:  config.DateEqual,
					Comparison: "01.01.2025",
				},
			},
			wantCount:   1,
			description: "Only one transaction on 01.01.2025",
		},
		{
			name: "date filter - greater than",
			filters: []config.Filter{
				config.DateFilter{
					Field:      "date",
					Condition:  config.DateGreaterThan,
					Comparison: "01.01.2025",
				},
			},
			wantCount:   2,
			description: "Two transactions after 01.01.2025",
		},
		{
			name: "number filter - equal",
			filters: []config.Filter{
				config.NumberFilter{
					Field:      "value",
					Condition:  config.NumberEqual,
					Comparison: 200.0,
				},
			},
			wantCount:   1,
			description: "Only one transaction with value 200",
		},
		{
			name: "number filter - greater than",
			filters: []config.Filter{
				config.NumberFilter{
					Field:      "value",
					Condition:  config.NumberGreaterThan,
					Comparison: 100.0,
				},
			},
			wantCount:   2,
			description: "Two transactions with value > 100",
		},
		{
			name: "string filter - contain",
			filters: []config.Filter{
				config.StringFilter{
					Field:      "description",
					Condition:  config.StringContain,
					Comparison: "Test",
				},
			},
			wantCount:   1,
			description: "Only one transaction containing 'Test'",
		},
		{
			name: "multiple filters - AND logic",
			filters: []config.Filter{
				config.DateFilter{
					Field:      "date",
					Condition:  config.DateGreaterThan,
					Comparison: "01.01.2025",
				},
				config.NumberFilter{
					Field:      "value",
					Condition:  config.NumberLessThan,
					Comparison: 300.0,
				},
			},
			wantCount:   1,
			description: "One transaction matching both filters (date > 01.01.2025 AND value < 300)",
		},
		{
			name: "filter that matches none",
			filters: []config.Filter{
				config.NumberFilter{
					Field:      "value",
					Condition:  config.NumberEqual,
					Comparison: 999.0,
				},
			},
			wantCount:   0,
			description: "No transactions match value 999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapters.FilterTransactions(mockTxs, tt.filters)

			if len(result) != tt.wantCount {
				t.Errorf("FilterTransactions() returned %d items, want %d. %s", len(result), tt.wantCount, tt.description)
			}
		})
	}
}

func TestFilterTransactions_WithSwedbankTransactions(t *testing.T) {
	// Test with actual Swedbank transactions
	txs := []adapters.SwedbankTransaction{
		{
			AccountNumber:   "LV02HABA0123456789012",
			EntryType:       adapters.SwedbankEntryTransaction,
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AccountHolder:   "TEST USER",
			Description:     "PAYMENT",
			Value:           100,
			Currency:        "EUR",
			Flow:            adapters.SwedbankCredit,
			ArchiveCode:     "2025010101234567",
			TransactionType: adapters.SwedbankTransactionToBank,
			ReferenceNumber: "",
			DocumentNumber:  "",
		},
		{
			AccountNumber:   "LV02HABA0123456789012",
			EntryType:       adapters.SwedbankEntryStartBalance,
			Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AccountHolder:   "",
			Description:     "",
			Value:           1000,
			Currency:        "EUR",
			Flow:            adapters.SwedbankCredit,
			ArchiveCode:     "",
			TransactionType: adapters.SwedbankTransactionStartBalance,
			ReferenceNumber: "",
			DocumentNumber:  "",
		},
	}

	// Filter to only include actual transactions (EntryType = "20")
	filters := []config.Filter{
		config.StringFilter{
			Field:      "Ieraksta tips",
			Condition:  config.StringEqual,
			Comparison: "20",
		},
	}

	adapted := adapters.AdaptTransactions(txs)
	result := adapters.FilterTransactions(adapted, filters)

	if len(result) != 1 {
		t.Errorf("FilterTransactions() returned %d items, want 1", len(result))
	}

	// Verify it's the transaction and not the balance
	if len(result) > 0 {
		normalized := result[0].Normalize()
		if normalized.Description != "PAYMENT" {
			t.Errorf("Filtered wrong transaction, got description %q", normalized.Description)
		}
	}
}
