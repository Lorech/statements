package config

import (
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const DefaultConfig = "config.json"
const schema = "schema/config.schema.json"

// Validates a configuration file against the JSON schema.
//
// If config is an empty string, the default config file `config.json` will be used.
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
