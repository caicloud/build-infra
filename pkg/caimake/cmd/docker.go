package cmd

import (
	"os"

	"github.com/caicloud/nirvana/cli"
	"github.com/spf13/cobra"
)

// NewCmdDocker ...
func NewCmdDocker() *cli.Command {
	cmd := cli.NewCommand(&cobra.Command{
		Use:   "docker",
		Short: "make for docker",
		Run:   runHelp,
	})
	cmd.AddCommand(
		newCmdDockerBuild(),
		newCmdDockerPush(),
	)
	return cmd
}

func newCmdDockerBuild() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "build",
		Short: "docker build",
		RunE: func(cmd *cobra.Command, args []string) error {
			code, err := bash.Run2("docker", "build", args...)
			if code > 0 {
				os.Exit(code)
			}
			return err
		},
	})
}

func newCmdDockerPush() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "push",
		Short: "docker push",
		RunE: func(cmd *cobra.Command, args []string) error {
			code, err := bash.Run2("docker", "push", args...)
			if code > 0 {
				os.Exit(code)
			}
			return err
		},
	})
}
