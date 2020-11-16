package main

import (
	"os"

	"github.com/athea-net/peggy/module/cmd/peggy/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
