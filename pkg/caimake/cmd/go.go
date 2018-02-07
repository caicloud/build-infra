package cmd

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			bash.Run2("go", "build", args...)
		},
	})
}

func newCmdGolangUnittest() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "unittest",
		Short: "golang unittest",
		Run: func(cmd *cobra.Command, args []string) {
			bash.Run2("go", "unittest", args...)
		},
	})
}
