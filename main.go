package main

import (
	"os"

	"github.com/visheyra/pbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
