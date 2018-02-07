package cmd

import (
	"github.com/caicloud/nirvana/cli"
	"github.com/spf13/cobra"
)

// NewCmdClean ...
func NewCmdClean() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "clean",
		Short: "clean make outputs",
		Run: func(cmd *cobra.Command, args []string) {
			bash.Run1("clean", args...)
		},
	})
}
