package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"statements/pkg/config"
	"testing"
)

func TestValidate(t *testing.T) {
	// Get the repository root
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Skipf("Skipping Validate tests: %v", err)
	}

	tests := []struct {
		name      string
		config    string
		content   string
		wantErr   bool
		errContain string
	}{
		{
			name:   "valid config",
			config: "test_valid_config.json",
			content: `{
				"$schema": "./schema/config.schema.json",
				"flags": {
					"bank": "swedbank",
					"output": "out.csv"
				},
				"filters": [
					{
						"field": "Ieraksta tips",
						"condition": "EQUAL",
						"comparison": "20"
					}
				]
			}`,
			wantErr: false,
		},
		{
			name:   "valid config with minimal flags",
			config: "test_minimal_config.json",
			content: `{
				"$schema": "./schema/config.schema.json",
				"flags": {
					"bank": "swedbank"
				},
				"filters": []
			}`,
			wantErr: false,
		},
		{
			name:       "invalid config - missing required fields",
			config:     "test_invalid_config.json",
			content:    `{"flags": {}}`,
			wantErr:    true,
			errContain: "invalid",
		},
		{
			name:       "invalid config - malformed JSON",
			config:     "test_malformed_config.json",
			content:    `{"flags": {`,
			wantErr:    true,
			errContain: "parse",
		},
		{
			name:       "nonexistent config file",
			config:     "nonexistent_file.json",
			content:    "",
			wantErr:    true,
			errContain: "could not open",
		},
	}

	// Save current directory and change to repo root for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	err = os.Chdir(repoRoot)
	if err != nil {
		t.Fatalf("Failed to change to repository root: %v", err)
	}

	// Create a temporary directory for test files
	tmpDir := t.TempDir()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip creating file for nonexistent file test
			var configPath string
			if tt.content != "" {
				configPath = filepath.Join(tmpDir, tt.config)
				err := os.WriteFile(configPath, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test config file: %v", err)
				}
			} else {
				configPath = filepath.Join(tmpDir, tt.config)
			}

			err := config.Validate(configPath)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
				} else if tt.errContain != "" && !containsString(err.Error(), tt.errContain) {
					t.Errorf("Validate() error = %v, should contain %q", err, tt.errContain)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		content    string
		wantBank   string
		wantOutput string
		wantErr    bool
		errContain string
	}{
		{
			name:   "valid config with all flags",
			config: "test_parse_full.json",
			content: `{
				"flags": {
					"bank": "swedbank",
					"input": "input.csv",
					"output": "output.csv"
				},
				"filters": []
			}`,
			wantBank:   "swedbank",
			wantOutput: "output.csv",
			wantErr:    false,
		},
		{
			name:   "valid config with partial flags",
			config: "test_parse_partial.json",
			content: `{
				"flags": {
					"bank": "swedbank"
				},
				"filters": []
			}`,
			wantBank:   "swedbank",
			wantOutput: "",
			wantErr:    false,
		},
		{
			name:       "invalid JSON",
			config:     "test_parse_invalid.json",
			content:    `{"flags": {`,
			wantErr:    true,
			errContain: "parse",
		},
		{
			name:       "nonexistent file",
			config:     "nonexistent.json",
			content:    "",
			wantErr:    true,
			errContain: "could not open",
		},
	}

	tmpDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var configPath string
			if tt.content != "" {
				configPath = filepath.Join(tmpDir, tt.config)
				err := os.WriteFile(configPath, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test config file: %v", err)
				}
			} else {
				configPath = filepath.Join(tmpDir, tt.config)
			}

			cfg, err := config.Parse(configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error but got none")
				} else if tt.errContain != "" && !containsString(err.Error(), tt.errContain) {
					t.Errorf("Parse() error = %v, should contain %q", err, tt.errContain)
				}
			} else {
				if err != nil {
					t.Errorf("Parse() unexpected error: %v", err)
					return
				}
				if cfg.Flags.Bank != tt.wantBank {
					t.Errorf("Parse() bank = %v, want %v", cfg.Flags.Bank, tt.wantBank)
				}
				if cfg.Flags.Output != tt.wantOutput {
					t.Errorf("Parse() output = %v, want %v", cfg.Flags.Output, tt.wantOutput)
				}
			}
		})
	}
}

func TestParseDefaultConfig(t *testing.T) {
	// Test that Parse uses DefaultConfig when empty string is provided
	tmpDir := t.TempDir()
	
	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create a config.json in the temp directory
	content := `{
		"flags": {
			"bank": "swedbank"
		},
		"filters": []
	}`
	err = os.WriteFile(config.DefaultConfig, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create default config file: %v", err)
	}

	cfg, err := config.Parse("")
	if err != nil {
		t.Errorf("Parse(\"\") unexpected error: %v", err)
		return
	}

	if cfg.Flags.Bank != "swedbank" {
		t.Errorf("Parse(\"\") bank = %v, want swedbank", cfg.Flags.Bank)
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Helper function to find the repository root
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repository root (no go.mod found)")
		}
		dir = parent
	}
}
