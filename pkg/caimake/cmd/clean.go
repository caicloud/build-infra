package cmd

import (
	"os"

	"github.com/caicloud/nirvana/cli"
	"github.com/spf13/cobra"
)

// NewCmdClean ...
func NewCmdClean() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "clean",
		Short: "clean make outputs",
		RunE: func(cmd *cobra.Command, args []string) error {
			code, err := bash.Run1("clean", args...)
			if code > 0 {
				os.Exit(code)
			}
			return err
		},
	})
}
