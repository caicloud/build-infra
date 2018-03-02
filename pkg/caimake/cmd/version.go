package cmd

import (
	"fmt"

	"github.com/caicloud/build-infra/pkg/version"
	"github.com/caicloud/nirvana/cli"
	"github.com/spf13/cobra"
)

// NewCmdVersion ...
func NewCmdVersion() *cli.Command {
	cmd := cli.NewCommand(&cobra.Command{
		Use:   "version",
		Short: "show caimake version information",
		Run: func(cmd *cobra.Command, args []string) {
			if cli.GetBool("short") {
				fmt.Println(version.Get().Version)
			} else {
				fmt.Println(version.Get().Pretty())
			}
		},
	})

	cmd.AddFlag(cli.BoolFlag{
		Name:      "short",
		Shorthand: "s",
		Usage:     "Show caimake version concisely",
	})
	return cmd
}
