package main

import (
	"os"

	"github.com/mimatache/forge/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		os.Exit(1)
	}
}
