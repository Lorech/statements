package main

import (
	"fmt"
	"os"
	"statements/pkg/commands"
)

func main() {
	cmd := commands.NewRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
