package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Config struct {
	Flags   FlagConfig  `json:"flags"`
	Filters []RawFilter `json:"filters"`
}

const DefaultConfig = "config.json"
const schema = "schema/config.schema.json"

// Validates a configuration file against the JSON schema.
//
// If `config` is an empty string, the default config file `config.json` will be used.
func Validate(config string) error {
	c := jsonschema.NewCompiler()
	sch, err := c.Compile(schema)
	if err != nil {
		return fmt.Errorf("could not compile JSON schema: %v", err)
	}

	if config == "" {
		config = DefaultConfig
	}

	f, err := os.Open(config)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer f.Close()

	inst, err := jsonschema.UnmarshalJSON(f)
	if err != nil {
		return fmt.Errorf("could not parse config file: %v", err)
	}

	err = sch.Validate(inst)
	if err != nil {
		return fmt.Errorf("config file invalid: %v", err)
	}

	return nil
}

// Parses a configuration file into a struct.
//
// If `config` is an empty string, the default config file `config.json` will be used.
func Parse(config string) (Config, error) {
	var c Config

	if config == "" {
		config = DefaultConfig
	}

	f, err := os.Open(config)
	if err != nil {
		return c, fmt.Errorf("could not open config file: %v", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&c); err != nil {
		return c, fmt.Errorf("could not parse config file: %v", err)
	}

	return c, nil
}
