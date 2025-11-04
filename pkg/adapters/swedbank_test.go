package adapters_test

import (
	"statements/pkg/adapters"
	"testing"
	"time"
)

func TestNewSwedbankTransaction(t *testing.T) {
	tests := []struct {
		name    string
		row     []string
		want    adapters.SwedbankTransaction
		wantErr bool
	}{
		{"parse transaction",
			[]string{
				"LV02HABA0123456789012",
				"20",
				"31.10.2025",
				"TEST USER",
				"SOME DESCRIPTION",
				"0,83",
				"EUR",
				"K",
				"2025103101234567",
				"INB",
				"",
				"",
			}, adapters.SwedbankTransaction{
				AccountNumber:   "LV02HABA0123456789012",
				EntryType:       adapters.SwedbankEntryTransaction,
				Date:            time.Date(2025, time.Month(10), 31, 0, 0, 0, 0, time.UTC),
				AccountHolder:   "TEST USER",
				Description:     "SOME DESCRIPTION",
				Value:           83,
				Currency:        "EUR",
				Flow:            adapters.SwedbankCredit,
				ArchiveCode:     "2025103101234567",
				TransactionType: adapters.SwedbankTransactionToBank,
				ReferenceNumber: "",
				DocumentNumber:  "",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := adapters.NewSwedbankTransaction(tt.row)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewSwedbankTransaction() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewSwedbankTransaction() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("NewSwedbankTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
