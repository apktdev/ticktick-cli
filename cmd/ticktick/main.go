package main

import (
	"fmt"
	"os"

	"github.com/apktdev/ticktick-cli/internal/cli"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func buildVersionString() string {
	return fmt.Sprintf("%s (commit %s, built %s)", version, commit, buildDate)
}

func main() {
	if err := cli.NewRootCmd(buildVersionString()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
