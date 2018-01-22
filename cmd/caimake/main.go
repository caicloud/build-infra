package main

import (
	"fmt"
	"os"

	"github.com/caicloud/build-infra/pkg/caimake/cmd"
)

func main() {
	cmd := cmd.NewCaimakeCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
