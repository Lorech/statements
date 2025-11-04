package commands

import (
	"os"
	"path/filepath"
	"statements/pkg/transactions"
	"testing"
	"time"
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantRows  int
		wantErr   bool
		errString string
	}{
		{
			name: "valid CSV with semicolon delimiter",
			content: `Header1;Header2;Header3
Value1;Value2;Value3
Value4;Value5;Value6`,
			wantRows:  3,
			wantErr:   false,
			errString: "",
		},
		{
			name: "empty CSV",
			content: ``,
			wantRows:  0,
			wantErr:   false,
			errString: "",
		},
		{
			name: "CSV with single row",
			content: `Header1;Header2`,
			wantRows:  1,
			wantErr:   false,
			errString: "",
		},
		{
			name:      "nonexistent file",
			content:   "",
			wantRows:  0,
			wantErr:   true,
			errString: "could not be opened",
		},
	}

	tmpDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			
			if tt.name == "nonexistent file" {
				filePath = filepath.Join(tmpDir, "nonexistent.csv")
			} else {
				filePath = filepath.Join(tmpDir, "test.csv")
				err := os.WriteFile(filePath, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			records, err := readInput(filePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("readInput() expected error but got none")
				} else if tt.errString != "" && !contains(err.Error(), tt.errString) {
					t.Errorf("readInput() error = %v, should contain %q", err, tt.errString)
				}
			} else {
				if err != nil {
					t.Errorf("readInput() unexpected error: %v", err)
					return
				}
				if len(records) != tt.wantRows {
					t.Errorf("readInput() returned %d rows, want %d", len(records), tt.wantRows)
				}
			}
		})
	}
}

func TestReadInput_ValidateContent(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.csv")
	
	content := `Header1;Header2;Header3
Value1;Value2;Value3`
	
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	records, err := readInput(filePath)
	if err != nil {
		t.Fatalf("readInput() unexpected error: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("readInput() returned %d rows, want 2", len(records))
	}

	// Validate first row
	if len(records[0]) != 3 {
		t.Errorf("First row has %d columns, want 3", len(records[0]))
	}
	if records[0][0] != "Header1" {
		t.Errorf("First row first column = %q, want Header1", records[0][0])
	}

	// Validate second row
	if len(records[1]) != 3 {
		t.Errorf("Second row has %d columns, want 3", len(records[1]))
	}
	if records[1][0] != "Value1" {
		t.Errorf("Second row first column = %q, want Value1", records[1][0])
	}
}

func TestWriteOutput(t *testing.T) {
	tests := []struct {
		name         string
		transactions []transactions.Transaction
		wantErr      bool
		errString    string
	}{
		{
			name: "write single transaction",
			transactions: []transactions.Transaction{
				{
					Date:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					AccountHolder: "Test User",
					Description:   "Test transaction",
					Value:         10050, // 100.50 EUR in cents
					Currency:      "EUR",
				},
			},
			wantErr:   false,
			errString: "",
		},
		{
			name: "write multiple transactions",
			transactions: []transactions.Transaction{
				{
					Date:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					AccountHolder: "User 1",
					Description:   "Transaction 1",
					Value:         10000,
					Currency:      "EUR",
				},
				{
					Date:          time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					AccountHolder: "User 2",
					Description:   "Transaction 2",
					Value:         20000,
					Currency:      "EUR",
				},
			},
			wantErr:   false,
			errString: "",
		},
		{
			name:         "write empty list",
			transactions: []transactions.Transaction{},
			wantErr:      false,
			errString:    "",
		},
	}

	tmpDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, "output.csv")

			err := writeOutput(filePath, tt.transactions)

			if tt.wantErr {
				if err == nil {
					t.Errorf("writeOutput() expected error but got none")
				} else if tt.errString != "" && !contains(err.Error(), tt.errString) {
					t.Errorf("writeOutput() error = %v, should contain %q", err, tt.errString)
				}
			} else {
				if err != nil {
					t.Errorf("writeOutput() unexpected error: %v", err)
					return
				}

				// Verify file was created
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("writeOutput() did not create output file")
				}
			}
		})
	}
}

func TestWriteOutput_ValidateContent(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "output.csv")

	txs := []transactions.Transaction{
		{
			Date:          time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			AccountHolder: "John Doe",
			Description:   "Salary payment",
			Value:         150000, // 1500.00 EUR in cents
			Currency:      "EUR",
		},
	}

	err := writeOutput(filePath, txs)
	if err != nil {
		t.Fatalf("writeOutput() unexpected error: %v", err)
	}

	// Read back the file and validate
	records, err := readInput(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("Output file has %d rows, want 1", len(records))
	}

	row := records[0]
	if len(row) != 5 {
		t.Errorf("Output row has %d columns, want 5", len(row))
	}

	// Validate the date
	if row[0] != "15.01.2025" {
		t.Errorf("Date = %q, want 15.01.2025", row[0])
	}

	// Validate account holder
	if row[1] != "John Doe" {
		t.Errorf("AccountHolder = %q, want John Doe", row[1])
	}

	// Validate description
	if row[2] != "Salary payment" {
		t.Errorf("Description = %q, want Salary payment", row[2])
	}

	// Validate value (should be formatted as 1500,0)
	if row[3] != "1500,0" {
		t.Errorf("Value = %q, want 1500,0", row[3])
	}

	// Validate currency
	if row[4] != "EUR" {
		t.Errorf("Currency = %q, want EUR", row[4])
	}
}

func TestWriteOutput_InvalidDirectory(t *testing.T) {
	// Try to write to a directory that doesn't exist and can't be created
	filePath := "/nonexistent_dir_12345/output.csv"

	txs := []transactions.Transaction{
		{
			Date:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AccountHolder: "Test",
			Description:   "Test",
			Value:         100,
			Currency:      "EUR",
		},
	}

	err := writeOutput(filePath, txs)
	if err == nil {
		t.Error("writeOutput() expected error for invalid directory but got none")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
