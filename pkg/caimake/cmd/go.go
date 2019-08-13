package cmd

import (
	"os"

	"github.com/caicloud/nirvana/cli"
	"github.com/spf13/cobra"
)

// NewCmdGolang ...
func NewCmdGolang() *cli.Command {
	cmd := cli.NewCommand(&cobra.Command{
		Use:   "go",
		Short: "make for golang",
		Run:   runHelp,
	})
	cmd.AddCommand(
		newCmdGolangBuild(),
		newCmdGolangUnittest(),
	)
	return cmd
}

func newCmdGolangBuild() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "build",
		Short: "golang build",
		RunE: func(cmd *cobra.Command, args []string) error {
			code, err := bash.Run2("go", "build", args...)
			if code > 0 {
				os.Exit(code)
			}
			return err
		},
	})
}

func newCmdGolangUnittest() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "unittest",
		Short: "golang unittest",
		RunE: func(cmd *cobra.Command, args []string) error {
			code, err := bash.Run2("go", "unittest", args...)
			if code > 0 {
				os.Exit(code)
			}
			return err
		},
	})
}
