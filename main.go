package main

import (
	"os"

	"github.com/cruious-kitten/forge/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		os.Exit(1)
	}
}
