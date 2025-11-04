package config

// FlagConfig defined in the configuration file.
type FlagConfig struct {
	Bank   string `json:"bank,omitempty"`
	Input  string `json:"input,omitempty"`
	Output string `json:"output,omitempty"`
}
